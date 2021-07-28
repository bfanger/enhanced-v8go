type IoWriter = {
  Write: (s: string) => [number, string | null];
}
declare interface go {
  "console": IoWriter;
}
