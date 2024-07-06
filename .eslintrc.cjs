module.exports = {
  root: true,
  env: { browser: true, es2020: true },
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:react-hooks/recommended",
  ],
  ignorePatterns: ["dist", ".eslintrc.cjs"],
  parser: "@typescript-eslint/parser",
  plugins: ["react-refresh"],
  rules: {
    "react-refresh/only-export-components": [
      "warn",
      { allowConstantExport: true },
    ],
  },
  extends: ["prettier"],
};

// module.exports = {
//   root: true,
//   env: {
//     browser: true,
//     es2020: true,
//   },
//   parser: "@typescript-eslint/parser",
//   parserOptions: {
//     ecmaVersion: 2020,
//     sourceType: "module",
//     project: "./tsconfig.json",
//   },
//   plugins: ["react-refresh", "@typescript-eslint"],
//   extends: [
//     "eslint:recommended",
//     "plugin:@typescript-eslint/recommended",
//     "plugin:react-hooks/recommended",
//     "plugin:prettier/recommended", // Ensure this is last to prevent conflicts
//   ],
//   ignorePatterns: ["dist", ".eslintrc.cjs"],
//   rules: {
//     "react-refresh/only-export-components": [
//       "warn",
//       { allowConstantExport: true },
//     ],

//     "@typescript-eslint/no-explicit-any": "warn", // Example of a TypeScript specific rule
//     // Add more custom rules here
//   },
//   settings: {
//     react: {
//       version: "detect", // Automatically detect the react version
//     },
//   },
// };
