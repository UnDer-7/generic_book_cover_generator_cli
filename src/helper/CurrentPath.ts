export function CurrentPath(): string {
    // @ts-ignore: pkg is added by pkg lib
    return (process.pkg) ? process.cwd() : __dirname;
}