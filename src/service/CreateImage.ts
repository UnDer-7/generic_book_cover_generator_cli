import {createCanvas, loadImage, registerFont, Image, Canvas} from 'canvas';
import fs from 'fs';

import {AppConfigContext} from "../types/AppConfigContext";

function loadFonts(appContext: AppConfigContext) {
    const {
        path: { font: fontPath},
        fontFamily: fontName,
    } = appContext;

    try {
        console.log('Registering fonts')
        registerFont(fontPath, {family: fontName});
    } catch (err: unknown) {
        console.error('Failed load fonts. Error:', err);
        throw err;
    }
}

async function getBGImage({ path: { backgroundImage: backgroundImagePath } }: AppConfigContext) {
    try {
        console.log('Starting to load background image')
        const image = await loadImage(backgroundImagePath);
        console.log('Successfully loaded background image')
        return image;
    } catch (err: unknown) {
        console.error('Could not load background image. Error: ', err);
        throw err;
    }
}

function drawImage(image: Image, fontName: string, text1: string, text2: string) {
   try {
       console.log('Starting to draw the image')
       const canvas = createCanvas(image.width, image.height);
       const ctx = canvas.getContext('2d');

       // Draw the image onto the canvas
       ctx.drawImage(image, 0, 0);

       // Define a function to find the maximum font size to fit the text width
       function getMaxFontSize(text: string) {
           let fontSize = 400; // Start with a large font size
           ctx.font = `${fontSize}px "${fontName}"`;
           while (ctx.measureText(text).width > canvas.width) {
               fontSize -= 1;
               ctx.font = `${fontSize}px "${fontName}"`;
           }
           return fontSize;
       }

       // Get maximum font sizes for both pieces of text
       let fontSize1 = getMaxFontSize(text1);
       let fontSize2 = getMaxFontSize(text2);

       // Calculate positions
       const text1Width = ctx.measureText(text1).width;
       const text2Width = ctx.measureText(text2).width;
       const x1 = (canvas.width - text1Width) / 2;
       const x2 = (canvas.width - text2Width) / 2;
       // Center 'Chapter' vertically and place '1' below
       const y1 = (canvas.height / 2) - fontSize2;
       const y2 = (canvas.height / 2) + fontSize1;

       // Set the fill style for the text and write both pieces of text onto the canvas
       ctx.fillStyle = 'white';
       ctx.font = `${fontSize1}px "${fontName}"`;
       ctx.fillText(text1, x1, y1);
       ctx.font = `${fontSize2}px "${fontName}"`;
       ctx.fillText(text2, x2, y2);

       console.log('Finished drawing the image')

       return canvas;
   } catch (err: unknown) {
       console.error('Could not draw the image. Err: ', err);
       throw err;
   }
}

function createFile(canvas: Canvas, text2: string, bookCoversOutputPath: string) {
    try {
        console.log('Starting to save the image to a file')
        const buffer = canvas.toBuffer('image/png');
        fs.writeFileSync(`${bookCoversOutputPath}/output_${text2}.png`, buffer);
        console.log('Image created. Out put path: ', bookCoversOutputPath);
    } catch (err: unknown) {
        console.error('Fail to create file. Error: ', err);
        throw err;
    }
}

export async function CreateImage(appContext: AppConfigContext, text1: string, text2: string): Promise<void> {
    const {
        path: {
            bookCoversOutput: bookCoversOutputPath
        },
        fontFamily: fontName,
    } = appContext;

    try {
        loadFonts(appContext);
        const image = await getBGImage(appContext);
        const canvas = drawImage(image, fontName, text1, text2);
        createFile(canvas, text2, bookCoversOutputPath);
    } catch (err: unknown) {
        console.error('Failed create image. Error:', err);
        throw err;
    }
}