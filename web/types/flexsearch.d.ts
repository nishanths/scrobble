//
// Copied from node_modules/flexsearch/index.d.ts, and pruned/modified.
//

declare module "flexsearch" {
  interface Index<T> {
    //TODO: Chaining
    readonly id: string;
    readonly index: string;
    readonly length: number;

    add(o: T): unknown;

    // Result without pagination -> T[]
    search(query: string, options: number | SearchOptions, callback: (results: T[]) => void): void;
    search(query: string, options?: number | SearchOptions): Promise<T[]>;
    search(options: SearchOptions & {query: string}, callback: (results: T[]) => void): void;
    search(options: SearchOptions & {query: string}): Promise<T[]>;

    // Result with pagination -> SearchResults<T>
    search(query: string, options: number | SearchOptions & { page?: boolean | Cursor}, callback: (results: SearchResults<T>) => void): void;
    search(query: string, options?: number | SearchOptions & { page?: boolean | Cursor}): Promise<SearchResults<T>>;
    search(options: SearchOptions & {query: string, page?: boolean | Cursor}, callback: (results: SearchResults<T>) => void): void;
    search(options: SearchOptions & {query: string, page?: boolean | Cursor}): Promise<SearchResults<T>>;

    clear(): unknown;
    destroy(): unknown;
  }

  interface SearchOptions {
    limit?: number,
    suggest?: boolean,
    where?: {[key: string]: string},
    field?: string | string[],
    bool?: "and" | "or" | "not"
    //TODO: Sorting
  }

  interface SearchResults<T> {
    page?: Cursor,
    next?: Cursor,
    result: T[]
  }

  interface Document {
      id: string;
      field: any;
  }


  export type CreateOptions = {
    profile?: IndexProfile;
    tokenize?: DefaultTokenizer | TokenizerFn;
    split?: RegExp;
    encode?: DefaultEncoder | EncoderFn | false;
    cache?: boolean | number;
    async?: boolean;
    worker?: false | number;
    depth?: false | number;
    threshold?: false | number;
    resolution?: number;
    stemmer?: Stemmer | string | false;
    filter?: FilterFn | string | false;
    rtl?: boolean;
    doc?: Document;
  };

  type IndexProfile = "memory" | "speed" | "match" | "score" | "balance" | "fast";
  type DefaultTokenizer = "strict" | "forward" | "reverse" | "full";
  type TokenizerFn = (str: string) => string[];
  type DefaultEncoder = "icase" | "simple" | "advanced" | "extra" | "balance";
  type EncoderFn = (str: string) => string;
  type Stemmer = {[key: string]: string};
  type Matcher = {[key: string]: string};
  type FilterFn = (str: string) => boolean;
  type Cursor = string;

  export default class FlexSearch {
    static create<T>(options?: CreateOptions): Index<T>;
  }
}
