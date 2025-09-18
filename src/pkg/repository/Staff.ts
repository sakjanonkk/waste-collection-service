import { Entity, PrimaryGeneratedColumn, Column, OneToMany } from "typeorm";
import { CollectionLog } from "./CollectionLog";
import { RouteReport } from "./RouteReport";
import { Vehicle } from "./Vehicle";

@Entity()
export class Staff {
  @PrimaryGeneratedColumn()
  s_id!: number;

  @Column({ type: "varchar", length: 30 })
  s_prename!: string;

  @Column({ type: "varchar", length: 50 })
  s_name!: string;

  @Column({ type: "varchar", length: 50 })
  s_surname!: string;

  @Column({ type: "enum", enum: ["MALE", "FEMALE", "OTHER"], nullable: true })
  s_gender!: string;

  @Column({ type: "varchar", length: 100, unique: true })
  s_email!: string;

  @Column({ type: "varchar", length: 255 })
  s_password!: string;

  @Column({ type: "enum", enum: ["PLANNER", "DRIVER"] })
  position!: string;

  @Column({ type: "enum", enum: ["ACTIVE", "INACTIVE"], default: "ACTIVE" })
  status!: string;

  @Column({ type: "varchar", length: 15, nullable: true })
  s_phone!: string;

  @Column({ type: "varchar", length: 255, nullable: true })
  s_image!: string;

  @Column({ type: "varchar", length: 6, nullable: true })
  otp_code!: string;

  @Column({ type: "timestamp", nullable: true })
  otp_expiry!: Date;

  @OneToMany(() => Vehicle, (vehicle) => vehicle.currentDriver)
  vehiclesDriven!: Vehicle[];

  @OneToMany(() => RouteReport, (routeReport) => routeReport.driver)
  routeReports!: RouteReport[];

  @OneToMany(() => CollectionLog, (collectionLog) => collectionLog.collectedBy)
  collectionLogs!: CollectionLog[];
}
