import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  ManyToOne,
  JoinColumn,
  OneToMany,
} from "typeorm";
import { Staff } from "./Staff";
import { Vehicle } from "./Vehicle";
import { CollectionLog } from "./CollectionLog";

@Entity()
export class RouteReport {
  @PrimaryGeneratedColumn()
  route_id!: number;

  @Column({ type: "date" })
  route_date!: Date;

  @Column()
  driver_id!: number;

  @Column()
  vehicle_id!: number;

  @Column({ type: "jsonb" })
  point_sequence: any;

  @Column({ type: "decimal", precision: 10, scale: 2, nullable: true })
  estimated_distance!: number;

  @Column({ type: "int", nullable: true })
  estimated_time!: number;

  @Column({
    type: "enum",
    enum: ["PENDING", "IN_PROGRESS", "COMPLETED"],
    default: "PENDING",
  })
  status!: string;

  @ManyToOne(() => Staff, (staff) => staff.routeReports)
  @JoinColumn({ name: "driver_id" })
  driver!: Staff;

  @ManyToOne(() => Vehicle, (vehicle) => vehicle.routeReports)
  @JoinColumn({ name: "vehicle_id" })
  vehicle!: Vehicle;

  @OneToMany(() => CollectionLog, (collectionLog) => collectionLog.route)
  collectionLogs!: CollectionLog[];
}
