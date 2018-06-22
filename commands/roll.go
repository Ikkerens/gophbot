package commands

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ikkerens/gophbot"
	commands "github.com/ikkerens/gophbot/handlers"
	"go.uber.org/zap"
)

var diceRegExp = regexp.MustCompile("(?:([0-9]+?)?d([0-9]+))")

func init() {
	commands.AddCommand("roll", roll)
}

func roll(_ *discordgo.Session, cmd *commands.InvokedCommand) {
	var (
		query               = strings.Join(cmd.Args, " ")
		rawQuery            = query
		diceCount, diceSize int64
		err                 error
	)
	matches := diceRegExp.FindAllStringSubmatch(rawQuery, -1)

	for _, match := range matches {
		if match[1] == "" {
			match[1] = "1"
		}

		diceCount, err = strconv.ParseInt(match[1], 10, 32)
		if err != nil {
			break
		}
		diceSize, err = strconv.ParseInt(match[2], 10, 32)
		if err != nil {
			break
		}

		if diceSize == 0 {
			cmd.Reply("... d0? What do you expect me to do with that?")
			return
		}

		rolls := make([]string, diceCount)
		for i := range rolls {
			rolls[i] = strconv.FormatInt(rand.Int63n(diceSize)+1, 10)
		}

		query = strings.Replace(query, match[0], fmt.Sprintf("(%s)", strings.Join(rolls, "+")), -1)
	}

	var sum int
	if err == nil {
		sum, err = eval(query)
	}

	if err != nil {
		gophbot.Log.Error("Could not parse request", zap.Error(err))
		cmd.Reply("You're not making sense.")
		return
	}

	cmd.Reply(fmt.Sprintf("`%s` = %s = %d", rawQuery, query, sum))
}

func eval(expr string) (int, error) {
	astExpr, err := parser.ParseExpr(expr)
	if err != nil {
		return 0, err
	}

	return evalAstExpr(astExpr)
}

func evalAstExpr(exp ast.Expr) (int, error) {
	switch exp := exp.(type) {
	case *ast.ParenExpr:
		return evalAstExpr(exp.X)
	case *ast.BinaryExpr:
		return evalBinaryExpr(exp)
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT:
			i, _ := strconv.Atoi(exp.Value)
			return i, nil
		default:
			return 0, errors.New("Unknown kind: " + string(exp.Kind))
		}
	default:
		return 0, errors.New("Invalid expression: " + fmt.Sprintf("%T", exp))
	}
}

func evalBinaryExpr(exp *ast.BinaryExpr) (int, error) {
	left, err := evalAstExpr(exp.X)
	if err != nil {
		return 0, err
	}
	right, err := evalAstExpr(exp.Y)
	if err != nil {
		return 0, err
	}

	switch exp.Op {
	case token.ADD:
		return left + right, nil
	case token.SUB:
		return left - right, nil
	case token.MUL:
		return left * right, nil
	case token.QUO:
		if right == 0 {
			return 0, errors.New("division by 0")
		}

		return int(math.Ceil(float64(left) / float64(right))), nil
	default:
		return 0, errors.New("Unknown token: " + string(exp.Op))
	}
}
