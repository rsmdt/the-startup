
// This should pass - proper type annotations
function processValue(value: string | number): string {
  return value.toString();
}

interface User {
  name: string;
  age: number;
}

function greetUser(user: User): string {
  return `Hello, ${user.name}!`;
}
