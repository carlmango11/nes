declare global {
  export interface Window {
    Go: any;
    getDisplay: () => number[];
  }
}

export {};
