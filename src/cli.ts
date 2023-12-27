#! /usr/bin/env node

import fs from "fs";

import { getUserInfo } from "./entrypoint/UserInput";
import {AppConfigContext} from "./types/AppConfigContext";
import {DeepFreeze} from "./helper/DeepFreeze";
import {DeepReadonly} from "ts-essentials";
import {RunAsync} from "./helper/Run";
import {CreateImage} from "./service/CreateImage";
import * as path from "path";

function checkIfNecessariesFilesExists(appCtx: AppConfigContext): void {
    const {
        path: {
            bookFolder,
            bookCoversOutput,
            font,
            backgroundImage
        }
    } = appCtx;

    // if (!fs.existsSync(bookFolder)) {
    //     throw new Error(`Book folder (-f) path does not exists. Informed path: ${bookFolder}`);
    // }

    if (!fs.existsSync(bookCoversOutput)) {
        const msg = `Book cover output (-o) path does not exists. Informed path: ${bookCoversOutput}`;
        throw new Error(msg)
    }

    if (!fs.existsSync(font)) {
        const msg = `font path does not exists. This is an internal file, this should not happen. Informed path: ${font}`;
        throw new Error(msg)
    }

    if (!fs.existsSync(backgroundImage)) {
        const msg = `background image path does not exists. This is an internal file, this should not happen. Informed path: ${backgroundImage}`;
        throw new Error(msg)
    }

    console.log('All paths are valid, proceeding with the script')
}

function createAppConfig(): DeepReadonly<AppConfigContext> {
    const {
        bookPath: bookFolder,
        outPath: bookCoversOutput,
    } = getUserInfo();

    const appConfig: AppConfigContext = {
        fontFamily: 'Merriweather',
        path: {
            bookFolder,
            bookCoversOutput,
            font: path.join(__dirname, '/resources/font/Merriweather-Black.ttf'),
            backgroundImage: path.join(__dirname, '/resources/background/black_background.jpg'),
        }
    }

    return DeepFreeze(appConfig);
}

RunAsync(async () => {
    console.log('Starting script');

    const appConfig = createAppConfig();
    checkIfNecessariesFilesExists(appConfig);

    // ToDo: PKG not working
    try {
        await CreateImage(appConfig, 'Chapter', '2');
    } catch (err: unknown) {
        console.error('Error: ', err);
    }

    console.log('App Config: ', appConfig);
});