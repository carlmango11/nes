import {useEffect} from 'react'
import "./wasm_exec.js";

const FPS = 20;
const HEIGHT = 200;
const WIDTH = 240;
const PIXEL_SIZE = 3;

function initWasm() {
    const go = new window.Go();

    WebAssembly.instantiateStreaming(fetch("nes.wasm"), go.importObject).then(
        (result) => {
            go.run(result.instance);
        }
    );

    startDisplay();
}

function startDisplay() {
    const canvas = document.getElementById("display") as HTMLCanvasElement;
    const ctx = canvas.getContext("2d");
    if (!ctx) {
        return;
    }

    setInterval(() => renderDisplay(ctx), 1000 / FPS);
}

function renderDisplay(ctx: CanvasRenderingContext2D) {
    ctx.clearRect(0, 0, WIDTH, HEIGHT);

    const state = window.getDisplay();
    // console.log(state);

    for (let y = 0; y < HEIGHT; y++) {
        for (let x = 0; x < WIDTH; x++) {
            const i = x + (y * WIDTH);
            const v: number = state[i];
            if (v !== 0) {
                console.log("OMG " + v);
            }

            ctx.fillStyle = getColour(v);
            ctx.fillRect(x*PIXEL_SIZE, y*PIXEL_SIZE, 1+PIXEL_SIZE, 1+PIXEL_SIZE);
        }
    }
}

function getColour(v: number): string {
    switch (v) {
        case 0:
            return "rgb(255, 255, 100)";
        case 1:
            return "rgb(255, 0, 0)";
        case 2:
            return "rgb(0, 255, 0)";
        case 3:
            return "rgb(0, 0, 255)";
    }

    return "rgb(0, 0, 0)";
}

function App() {
    useEffect(() => {
        initWasm();
    }, []);

    return <div>
        <canvas id="display" width="720" height="600"></canvas>
    </div>
}

export default App
