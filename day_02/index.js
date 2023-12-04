import readline from 'readline';
import { TLSSocket } from 'tls';

const Game = 'Game';
const Blue = 'Blue';
const Green = 'Green';
const Red = 'Red';
const Colon = 'Colon';
const Semicolon = 'Semicolon';
const Comma = 'Comma';
const Int = 'Int';
const EOL = 'EOL';

/**
 * @readonly
 * @enum {string}
 **/
const TokenType = {
    Game,
    Blue,
    Green,
    Red,
    Colon,
    Semicolon,
    Comma,
    Int,
    EOL,
}

class Token {
    /**
     * @param type {TokenType}
     * @param literal {string}
     **/
    constructor(type, literal) {
        this.type = type;
        this.literal = literal;
    }
}

/**
 * @param char {string}
 **/
const isLetter = char => {
    const charCode = char.charCodeAt(0);
    // console.log({char, charCode});
    return (charCode >= 65 && charCode <= 90) || (charCode >= 97 && charCode <= 122);
}

/**
 * @param char {string}
 **/
const isDigit = char => {
    const charCode = char.charCodeAt(0);
    return charCode >= 48 && charCode <= 57;
}

/**
 * @param identifier {string}
 **/
const lookupTokenType = identifier => {
    switch(identifier) {
        case 'Game': return TokenType.Game;
        case 'blue': return TokenType.Blue;
        case 'red': return TokenType.Red;
        case 'green': return TokenType.Green;
    }
    return null;
}

class Lexer {

    /**
     * @param input {string}
     **/
    constructor(input) {
        this.input = input;
        this.position = 0;
        this.readPosition = 0;
        this.char = '';
        this.readChar();
    }

    readChar() {
        if (this.readPosition >= this.input.length) {
            this.char = '';
        } else {
            this.char = this.input[this.readPosition];
        }
        this.position = this.readPosition;
        this.readPosition += 1;
    }

    nextToken() {
        this.skipWhitespace();

        const char = this.char;

        switch(this.char) {
            case ':':
                this.readChar();
                return new Token(TokenType.Colon, char);
            case ';':
                this.readChar();
                return new Token(TokenType.Semicolon, char);
            case ',':
                this.readChar();
                return new Token(TokenType.Comma, char);
            case '':
                this.readChar();
                return new Token(TokenType.EOL, char);
            default:
                if (isLetter(this.char)) {
                    const identifier = this.readIdentifier();
                    const tokenType = lookupTokenType(identifier);
                    if (tokenType !== null) {
                        return new Token(tokenType, identifier);
                    } else {
                        throw new Error(`Unknown identifier "${identifier}"`);
                    }
                } else if (isDigit(this.char)) {
                    const number = this.readNumber();
                    return new Token(TokenType.Int, number);
                } else {
                    throw new Error(`Unknown char "${char}"`);
                }
        }
    }

    readNumber() {
        const position = this.position;
        while (isDigit(this.char)) {
            this.readChar();
        }
        return this.input.substring(position, this.position);
    }

    readIdentifier() {
        const position = this.position;
        while (isLetter(this.char)) {
            this.readChar();
        }
        return this.input.substring(position, this.position);
    }

    skipWhitespace() {
        while (this.char === ' ') {
            this.readChar();
        }
    }
}

class Statement {}

class GameStatement extends Statement {
    /**
     * @param gameId {number}
     * @param subsets {SubsetStatement[]}
     **/
    constructor(gameId, subsets) {
        super();
        this.gameId = gameId;
        this.subsets = subsets;
    }
}

class RevealStatement extends Statement {
    /**
     * @param color {'blue' | 'green' | 'red'}
     * @param amount {number}
     **/
    constructor(color, amount) {
        super();
        this.color = color;
        this.amount = amount;
    }
}

class SubsetStatement extends Statement {
    /**
     * @param reveals {RevealStatement[]}
     **/
    constructor(reveals) {
        super();
        this.reveals = reveals;
    }
}

class Parser {
    /**
     * @param lexer {Lexer}
     **/
    constructor(lexer) {
        this.lexer = lexer;
        this.curToken = null;
        this.peekToken = null;
        this.nextToken();
        this.nextToken();
    }

    parseGame() {
        /**
         * @type {GameStatement[]}
         **/
        const statements = [];
        while (!this.curTokenIs(TokenType.EOL)) {
            statements.push(this.parseStatement());
        }
        return statements;
    }

    nextToken() {
        this.curToken = this.peekToken;
        this.peekToken = this.lexer.nextToken();
    }

    /**
     * @param tokenType {TokenType}
     **/
    curTokenIs(tokenType) {
        return this.curToken.type === tokenType;
    }

    /**
     * @param tokenType {TokenType}
     **/
    peekTokenIs(tokenType) {
        return this.peekToken.type === tokenType;
    }

    parseStatement() {
        if (this.curTokenIs(TokenType.Game)) {
            return this.parseGameStatement();
        }
        throw new Error(`unexpected statement: "${this.curToken.type}"`);
    }

    parseGameStatement() {
        this.expectPeek(TokenType.Int);
        const id = parseInt(this.curToken.literal);
        this.expectPeek(TokenType.Colon);
        const subsets = this.parseSubsets();
        return new GameStatement(id, subsets);
    }

    parseRevealStatement() {
        this.expectPeek(TokenType.Int);
        const amount = parseInt(this.curToken.literal);
        if (this.peekTokenIs(TokenType.Green) || this.peekTokenIs(TokenType.Red) || this.peekTokenIs(TokenType.Blue)) {
            this.nextToken();
        }
        const color = this.curToken.literal;
        this.nextToken();
        return new RevealStatement(color, amount);
    }

    parseSubset() {
        const reveals = [];
        while (!this.curTokenIs(TokenType.EOL)) {
            reveals.push(this.parseRevealStatement());
            if (this.curTokenIs(TokenType.Semicolon)) {
                break;
            }
        }
        return new SubsetStatement(reveals);
    }

    parseSubsets() {
        const subsets = [];
        while(!this.curTokenIs(TokenType.EOL)) {
            subsets.push(this.parseSubset());
        }
        return subsets;
    }

    /**
     * @param tokenType {TokenType}
     **/
    expectPeek(tokenType) {
        if (!this.peekTokenIs(tokenType)) {
            throw new Error(`expected token "${tokenType}", got "${this.peekToken.type}"`);
        }
        this.nextToken();
    }
}

/**
 * @param line {string}
 **/
const parseGame = line => {
    const lexer = new Lexer(line);
    const parser = new Parser(lexer);
    const [gameStatement] = parser.parseGame();
    const subsets = [];
    for (const subsetStatement of gameStatement.subsets) {
        const subset = {
            red: 0,
            blue: 0,
            green: 0,
        };
        for (const revealStatement of subsetStatement.reveals) {
            subset[revealStatement.color] += revealStatement.amount;
        }
        subsets.push(subset);
    }
    return {subsets, gameId: gameStatement.gameId}
}

/**
 * @param game {{gameId: number, subsets: {red: number, green: number, blue: number}[]}}
 * @param predicate {{red: number, green: number, blue: number}}
 **/
const isValidGame = (game, predicate) => {
    for (const subset of game.subsets) {
        if (subset.red > predicate.red || subset.blue > predicate.blue || subset.green > predicate.green) {
            return false;
        }
    }
    return true;
}

/**
 * @param subsets {{red: number, blue: number, green: number}[]}
 **/
const calculatePower = subsets => {
    let red = 0;
    let green = 0;
    let blue = 0;
    for (const subset of subsets) {
        red = Math.max(red, subset.red);
        green = Math.max(green, subset.green);
        blue = Math.max(blue, subset.blue);
    }
    return red * green * blue;
}

const main = async () => {
    const rl = readline.createInterface({
        input: process.stdin,
    });

    let sum = 0;
    let sumOfPowers = 0;
    const predicate = {
        red: 12,
        green: 13,
        blue: 14,
    };
    for await (const line of rl) {
        const game = parseGame(line);
        if (isValidGame(game, predicate)) {
            sum += game.gameId;
        }
        sumOfPowers += calculatePower(game.subsets);
    }
    console.log(sum);
    console.log(sumOfPowers);
}

main();
