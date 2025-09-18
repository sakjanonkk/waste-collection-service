import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  ManyToOne,
  JoinColumn,
} from "typeorm";
import { CollectionPoint } from "./CollectionPoint";

@Entity()
export class PickupRequest {
  @PrimaryGeneratedColumn()
  request_id!: number;

  @Column({ type: "enum", enum: ["REPORT_ISSUE"] })
  request_type!: string;

  @Column()
  point_id!: number;

  @Column({ type: "decimal", precision: 10, scale: 8 })
  latitude!: number;

  @Column({ type: "decimal", precision: 11, scale: 8 })
  longitude!: number;

  @Column({
    type: "enum",
    enum: ["REGULAR", "RECYCLE", "BOTH"],
    nullable: true,
  })
  waste_type!: string;

  @Column({ type: "text" })
  remarks!: string;

  @CreateDateColumn({ type: "timestamp with time zone" })
  request_datetime!: Date;

  @Column({ type: "enum", enum: ["PENDING", "COMPLETED"], default: "PENDING" })
  status!: string;

  @ManyToOne(
    () => CollectionPoint,
    (collectionPoint) => collectionPoint.pickupRequests
  )
  @JoinColumn({ name: "point_id" })
  collectionPoint!: CollectionPoint;
}
