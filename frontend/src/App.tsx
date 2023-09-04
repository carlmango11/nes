import {useEffect} from 'react'
import "./wasm_exec.js";

function initWasm() {
    const go = new window.Go();

    WebAssembly.instantiateStreaming(fetch("nes.wasm"), go.importObject).then(
        (result) => {
            go.run(result.instance);
        }
    );
}

function clicky() {
    const d = window.getDisplay();
    console.log(d);
}

function App() {
    useEffect(() => {
        initWasm();
    }, []);

    return <div>
        <p onClick={clicky} >OMG</p>
    </div>
}

export default App
