export function Run(func: () => void): void {
    func();
}

export async function RunAsync(func: () => Promise<void>): Promise<void> {
    await func();
}