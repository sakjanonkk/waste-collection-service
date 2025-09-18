import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  ManyToOne,
  JoinColumn,
} from "typeorm";
import { RouteReport } from "./RouteReport";
import { CollectionPoint } from "./CollectionPoint";
import { Staff } from "./Staff";

@Entity()
export class CollectionLog {
  @PrimaryGeneratedColumn()
  log_id!: number;

  @Column()
  route_id!: number;

  @Column()
  point_id!: number;

  @Column({
    type: "timestamp with time zone",
    default: () => "CURRENT_TIMESTAMP",
  })
  collected_at!: Date;

  @Column()
  collected_by_id!: number;

  @Column({ type: "decimal", precision: 10, scale: 2, nullable: true })
  regular_waste_amount!: number;

  @Column({ type: "decimal", precision: 10, scale: 2, nullable: true })
  recycle_waste_amount!: number;

  @Column({ type: "enum", enum: ["COMPLETED", "SKIPPED", "PROBLEM"] })
  status!: string;

  @Column({ type: "text", nullable: true })
  notes!: string;

  @ManyToOne(() => RouteReport, (routeReport) => routeReport.collectionLogs)
  @JoinColumn({ name: "route_id" })
  route!: RouteReport;

  @ManyToOne(
    () => CollectionPoint,
    (collectionPoint) => collectionPoint.collectionLogs
  )
  @JoinColumn({ name: "point_id" })
  collectionPoint!: CollectionPoint;

  @ManyToOne(() => Staff, (staff) => staff.collectionLogs)
  @JoinColumn({ name: "collected_by_id" })
  collectionLogs!: CollectionLog[];
  collectedBy: any;
}
