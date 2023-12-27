import yargs from 'yargs';
import fs from "fs";
import * as path from "path";
import {CurrentPath} from "../helper/CurrentPath";

export function getUserInfo(): {bookPath: string, outPath: string} {
    const defaultOutPath = `${CurrentPath()}/book_covers`;
    const argv = yargs(process.argv.slice(2))
        .options({
            f: { type: 'string', describe: 'Book folder', demandOption: true },
            o: { type: 'string', describe: 'Generated covers output path', default: defaultOutPath}
        }).parseSync();

    if (argv.o === defaultOutPath) {
        createDefaultOutPath(defaultOutPath)
    }

    return {
        bookPath: sanitizePath(argv.f),
        outPath: sanitizePath(argv.o),
    }
}

function createDefaultOutPath(path: string): void {
    const pathExits = fs.existsSync(path);
    console.log('does need to create default output path: ', pathExits)
    if (!path) {
        console.log('Starting to creating default output path. ', path)
        fs.mkdirSync(path);
        console.log('Finished creating default output path')
    }
}

function sanitizePath(path: string): string {
    if (path.startsWith("/")) {
        return path;
    }

    return '/' + path;
}