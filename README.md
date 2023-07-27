# kindsys2




typescript tools?


https://github.com/YousefED/typescript-json-schema

https://github.com/vega/ts-json-schema-generator


Composable:
1. Write typescript
2. generate jsonschema
```
npx ts-json-schema-generator --path example/transformers/filterByValue/options.ts --type 'Options' > example/transformers/filterByValue/options.schema.gen.json
```
3. generate golang
```
npx quicktype --package xxx --just-types -s schema example/transformers/filterByValue/options.schema.gen.json -o example/transformers/filterByValue/options.gen.go
```
