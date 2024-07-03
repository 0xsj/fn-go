// scripts/gen-json.ts

import * as fs from 'fs';
import * as path from 'path';
import { exec } from 'child_process';

// Define paths
const typesFilePath = path.resolve(__dirname, '..', 'src', 'types', 'types.ts');
const jsonSchemaPath = path.resolve(__dirname, '..', 'schema.json'); // Output JSON Schema path

// Command to generate JSON Schema
const command = `typescript-json-schema --refs --strictNullChecks ${typesFilePath} Employee Address --required -o ${jsonSchemaPath}`;

// Execute the command
exec(command, (error, stdout, stderr) => {
  if (error) {
    console.error(`Error executing command: ${error.message}`);
    return;
  }
  if (stderr) {
    console.error(`Command stderr: ${stderr}`);
    return;
  }
  console.log(`JSON Schema generated successfully at ${jsonSchemaPath}`);

  // Read generated JSON Schema file
  fs.readFile(jsonSchemaPath, 'utf8', (err, data) => {
    if (err) {
      console.error(`Error reading JSON Schema file: ${err.message}`);
      return;
    }
    const jsonSchema = JSON.parse(data);

    // Optionally, you can further process or use the JSON schema here
    console.log('Generated JSON Schema:', jsonSchema);
  });
});
