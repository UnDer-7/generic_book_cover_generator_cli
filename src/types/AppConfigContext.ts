export interface AppConfigContext {
    path: AppPathConfig
    fontFamily: string,
}

export interface AppPathConfig {
    bookFolder: string,
    backgroundImage: string,
    font: string,
    bookCoversOutput: string
}