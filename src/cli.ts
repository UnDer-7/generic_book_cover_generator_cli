#! /usr/bin/env node

import { getUserInfo } from "./entrypoint/UserInput";
import {AppConfig} from "./types/AppConfig";
import {DeepFreeze} from "./helper/DeepFreeze";
import {DeepReadonly} from "ts-essentials";
import {Run} from "./helper/Run";

function createAppConfig(): DeepReadonly<AppConfig> {
    const userInfo = getUserInfo();
    const appConfig: AppConfig = {
        fontFamily: 'Merriweather',
        path: {
            bookFolder: userInfo.bookPath,
            font: './resources/font/Merriweather-Black.ttf',
            backgroundImage: './resources/background/black_background.jpg'
        }
    }

    return DeepFreeze(appConfig);
}

Run(() => {
    console.log('Starting script');

    const appConfig = createAppConfig();

    console.log('App Config: ', appConfig);
});