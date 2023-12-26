#! /usr/bin/env node

import { getUserInfo } from "./entrypoint/UserInput";
import {AppConfigContext} from "./types/AppConfigContext";
import {DeepFreeze} from "./helper/DeepFreeze";
import {DeepReadonly} from "ts-essentials";
import {RunAsync} from "./helper/Run";
import {CreateImage} from "./service/CreateImage";

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
            font: `${__dirname}/resources/font/Merriweather-Black.ttf`,
            backgroundImage: `${__dirname}/resources/background/black_background.jpg`,
        }
    }

    return DeepFreeze(appConfig);
}

RunAsync(async () => {
    console.log('Starting script');

    const appConfig = createAppConfig();

    // ToDo: PKG not working
    await CreateImage(appConfig, 'Chapter', '2');

    console.log('App Config: ', appConfig);
});