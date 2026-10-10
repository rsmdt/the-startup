export interface User {
  id: string;
  email: string;
  displayName: string;
}

export interface UserRepository {
  findByEmail(email: string): Promise<User | null>;
  insert(user: Omit<User, "id">): Promise<User>;
}

export class UserService {
  constructor(private readonly users: UserRepository) {}

  async register(email: string, displayName: string): Promise<User> {
    const normalized = email.trim().toLowerCase();
    if (await this.users.findByEmail(normalized)) {
      throw new Error("email already registered");
    }
    return this.users.insert({ email: normalized, displayName });
  }
}
