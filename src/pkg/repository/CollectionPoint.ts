import { Entity, PrimaryGeneratedColumn, Column, OneToMany } from "typeorm";
import { PickupRequest } from "./PickupRequest";
import { CollectionLog } from "./CollectionLog";

@Entity()
export class CollectionPoint {
  @PrimaryGeneratedColumn()
  point_id!: number;

  @Column({ type: "varchar", length: 100 })
  point_name!: string;

  @Column({ type: "decimal", precision: 10, scale: 8 })
  latitude!: number;

  @Column({ type: "decimal", precision: 11, scale: 8 })
  longitude!: number;

  @Column({ type: "enum", enum: ["COLLECTION", "DEPOT"] })
  point_type!: string;

  @Column({ type: "enum", enum: ["ACTIVE", "INACTIVE"], default: "ACTIVE" })
  status!: string;

  @Column({ type: "varchar", length: 255, nullable: true })
  point_image!: string;

  @OneToMany(
    () => PickupRequest,
    (pickupRequest) => pickupRequest.collectionPoint
  )
  pickupRequests!: PickupRequest[];

  @OneToMany(
    () => CollectionLog,
    (collectionLog) => collectionLog.collectionPoint
  )
  collectionLogs: CollectionLog[] | undefined;
}
