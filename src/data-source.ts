import "reflect-metadata";
import { DataSource } from "typeorm";

export const AppDataSource = new DataSource({
  type: "postgres",
  host: "localhost",
  port: 5432,
  username: "postgres",
  password: "1234",
  database: "waste_management_db",
  synchronize: false,
  logging: true,
  entities: [__dirname + "/entity/*.ts"],
  migrations: [],
  subscribers: [],
});
