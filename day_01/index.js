import readline from 'readline';

const numbers = {
    1: 1,
    2: 2,
    3: 3,
    4: 4,
    5: 5,
    6: 6,
    7: 7,
    8: 8,
    9: 9,
}

const numberNames = {
    one: 1,
    two: 2,
    three: 3,
    four: 4,
    five: 5,
    six: 6,
    seven: 7,
    eight: 8,
    nine: 9,
}

const names = Object.keys(numberNames);

/**
 * @param line {string}
 * @param index {number}
 **/
const readNameAt = (line, index) => {
    for (const name of names) {
        const n = line.substring(index, index + name.length);
        if (n === name) {
            return name;
        }
    }
    return null;
}

/**
 * @param line {string}
 **/
const readFirstNum = line => {
    const len = line.length;
    for (let i = 0; i < len; i++) {
        const char = line[i];
        const num = numbers[char];
        if (typeof num === 'number') {
            return num;
        }
        const name = readNameAt(line, i);
        if (name !== null) {
            return numberNames[name] ?? 0;
        }
    }
    return 0;
}

/**
 * @param line {string}
 **/
const readLastNum = line => {
    let len = line.length;
    while (len) {
        const char = line[--len];
        const num = numbers[char];
        if (typeof num === 'number') {
            return num;
        }
        const name = readNameAt(line, len);
        if (name !== null) {
            return numberNames[name] ?? 0;
        }
    }
    return 0;
}

/**
 * @param line {string}
 **/
const readCalibration = line => {
    const firstNum = readFirstNum(line);
    const lastNum = readLastNum(line);
    return firstNum * 10 + lastNum;
}

const main = async () => {
    const rl = readline.createInterface({
        input: process.stdin,
    });

    let sum = 0;
    for await (const line of rl) {
        const calibration = readCalibration(line);
        sum += calibration;
    }
    console.log(sum);
}

main();

