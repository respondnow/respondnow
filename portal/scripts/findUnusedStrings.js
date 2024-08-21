const { execSync } = require('child_process');
const fs = require('fs');
const yaml = require('yaml');
const glob = require('glob');

// Function to load YAML file
function loadYAML(filePath) {
  try {
    const content = fs.readFileSync(filePath, 'utf8');
    return yaml.parse(content);
  } catch (error) {
    console.error(`Error reading YAML file: ${filePath}`);
    console.error(error);
    process.exit(1);
  }
}

// Function to find string occurrences in TypeScript and TypeScript with JSX files
function findStringOccurrences(yamlStrings, codebasePath, excludedFiles) {
  const stringKeys = Object.keys(flattenObject(yamlStrings));
  const notFoundStrings = new Set(stringKeys);

  // Create a regular expression to match the strings without quotes or backticks
  const stringRegex = new RegExp(`['"](${stringKeys.join(`|`)})['"]`, 'g');

  // Use glob to find TypeScript and TypeScript with JSX files in the codebase
  const tsFiles = glob.sync(`${codebasePath}/**/*.ts`);
  const tsxFiles = glob.sync(`${codebasePath}/**/*.tsx`);
  const allFiles = [...tsFiles, ...tsxFiles].filter(file => !excludedFiles.includes(file));

  // Iterate through each file and check for string occurrences
  allFiles.forEach(file => {
    const content = fs.readFileSync(file, 'utf8');
    const matches = content.match(stringRegex);

    if (matches) {
      matches.forEach(match => {
        const keyMatch = match.match(/['"]+(.+)['"]+/);

        if (keyMatch && keyMatch[1]) {
          const key = keyMatch[1];
          notFoundStrings.delete(key);
        }
      });
    }
  });

  return Array.from(notFoundStrings);
}

// Function to remove unused keys from the YAML file
function removeUnusedKeys(yamlFilePath, unusedKeys) {
  const existingYaml = loadYAML(yamlFilePath);

  unusedKeys.forEach(key => {
    deleteNestedKey(existingYaml, key);
  });

  const updatedYaml = yaml.stringify(existingYaml, { indent: 2 });

  fs.writeFileSync(yamlFilePath, updatedYaml, 'utf8');
}

// Recursive function to delete a nested key in an object
function deleteNestedKey(obj, key) {
  const keys = key.split('.');
  const lastKey = keys.pop();

  const parent = keys.reduce((acc, k) => acc && acc[k], obj);
  if (parent && parent.hasOwnProperty(lastKey)) {
    delete parent[lastKey];

    // Clean up empty nested objects
    keys.reduce((acc, k) => {
      if (acc && acc[k] && Object.keys(acc[k]).length === 0) {
        delete acc[k];
      }
      return acc && acc[k];
    }, obj);
  }
}

// Flatten a nested object
function flattenObject(obj, parentKey = '') {
  return Object.entries(obj).reduce((acc, [key, value]) => {
    const newKey = parentKey ? `${parentKey}.${key}` : key;

    if (typeof value === 'object' && value !== null) {
      Object.assign(acc, flattenObject(value, newKey));
    } else {
      acc[newKey] = value;
    }

    return acc;
  }, {});
}

// Function to process command line arguments
function processCommandLineArgs() {
  const args = process.argv.slice(2);

  if (args.includes('-r')) {
    return true;
  }

  return false;
}

// Run the yarn strings:sort command
function runYarnStringsSort() {
  try {
    execSync('yarn strings:sort', { stdio: 'inherit' });
  } catch (error) {
    console.error('Error running yarn strings:sort command.');
    console.error(error);
  }
}

// Specify the path to your strings YAML file, codebase, and excluded file
const codebasePath = 'src';
const stringsYamlFilePath = 'src/strings/strings.en.yaml';
const excludedFiles = ['src/strings/types.ts'];

// Load YAML file
const yamlStrings = loadYAML(stringsYamlFilePath);

// Find string occurrences in TypeScript and TypeScript with JSX files, excluding the specified file
const notFoundStrings = findStringOccurrences(yamlStrings, codebasePath, excludedFiles);

// Output not found strings
if (notFoundStrings.length > 0) {
  console.log('Strings not found in the codebase:');
  console.log(notFoundStrings.join(', '));

  // Process command line arguments
  const shouldRemoveUnusedKeys = processCommandLineArgs();

  // Remove unused keys from the YAML file if -r flag is provided
  if (shouldRemoveUnusedKeys) {
    removeUnusedKeys(stringsYamlFilePath, notFoundStrings);
    console.log('Unused keys removed from the YAML file.');

    // Run yarn strings:sort command
    console.log('Running yarn strings:sort command...');
    runYarnStringsSort();
  } else {
    console.log('\n\nTo remove unused keys, use the -r flag in the command line.');
  }
} else {
  runYarnStringsSort();
  console.log('All strings found in the codebase.');
}
