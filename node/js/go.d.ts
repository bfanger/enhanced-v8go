declare interface go {
  "os.Stdout": {
    WriteString: (s: string) => [int, string | null];
  };
}
