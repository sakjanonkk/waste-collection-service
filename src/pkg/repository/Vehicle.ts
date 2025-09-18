import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  ManyToOne,
  JoinColumn,
  OneToMany,
} from "typeorm";
import { Staff } from "./Staff";
import { RouteReport } from "./RouteReport";

@Entity()
export class Vehicle {
  @PrimaryGeneratedColumn()
  vehicle_id!: number;

  @Column({ type: "varchar", length: 20, unique: true })
  vehicle_reg_num!: string;

  @Column({
    type: "enum",
    enum: ["AVAILABLE", "IN_USE", "MAINTENANCE"],
    default: "AVAILABLE",
  })
  status!: string;

  @Column({ type: "int" })
  regular_capacity!: number;

  @Column({ type: "int" })
  recycle_capacity!: number;

  @Column({ nullable: true })
  current_driver_id!: number;

  @Column({ type: "text", nullable: true })
  problem_reported!: string;

  @ManyToOne(() => Staff, (staff) => staff.vehiclesDriven)
  @JoinColumn({ name: "current_driver_id" })
  currentDriver!: Staff;

  @OneToMany(() => RouteReport, (routeReport) => routeReport.vehicle)
  routeReports!: RouteReport[];
}
