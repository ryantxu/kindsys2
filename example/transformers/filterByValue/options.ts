// This is the source file

export enum FilterByValueType {
  exclude = "exclude",
  include = "include",
}

export enum FilterByValueMatch {
  all = "all",
  any = "any",
}

export interface FilterByValueFilter {
  fieldName: string;
  config: object; // MatcherConfig
}

export interface Options {
  filters: FilterByValueFilter[];
  type: FilterByValueType;
  match: FilterByValueMatch;
}
