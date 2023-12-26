export function DeepFreeze<T extends object>(obj: T): T {
    Object.keys(obj).forEach((prop: string) => {
        const objElement = obj[prop as keyof object];
        if (typeof objElement === 'object') {
            DeepFreeze(objElement);
        }
    });

    return Object.freeze(obj);
}