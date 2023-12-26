import yargs from 'yargs';

export function getUserInfo(): {bookPath: string} {
    const argv = yargs(process.argv.slice(2))
        .options({
            f: { type: 'string', describe: 'Book folder', demandOption: true }
        }).parseSync();

    return { bookPath: sanitizeBookPath(argv.f) }
}

function sanitizeBookPath(path: string): string {
    if (path.startsWith("/")) {
        return path;
    }

    return '/' + path;
}