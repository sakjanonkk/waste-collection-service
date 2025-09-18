import express from "express";
import { AppDataSource } from "./data-source";
import { Staff } from "./pkg/repository/Staff";

AppDataSource.initialize()
  .then(async () => {
    console.log("✅ Data Source has been initialized!");

    const app = express();
    app.use(express.json());

    app.get("/api/v1/staff", async (req, res) => {
      try {
        const staffRepository = AppDataSource.getRepository(Staff);
        const allStaff = await staffRepository.find();
        res.json(allStaff);
      } catch (error) {
        res.status(500).json({ message: "Error fetching staff", error });
      }
    });

    app.listen(3000, () => {
      console.log("🚀 Server started on http://localhost:3000");
    });
  })
  .catch((error) =>
    console.log("❌ Error during Data Source initialization", error)
  );
