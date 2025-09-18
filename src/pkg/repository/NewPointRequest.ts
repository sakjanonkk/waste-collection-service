import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
} from "typeorm";

@Entity()
export class NewPointRequest {
  @PrimaryGeneratedColumn()
  request_id!: number;

  @Column({ type: "enum", enum: ["NEW_POINT"] })
  request_type!: string;

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

  @Column({
    type: "enum",
    enum: ["PENDING", "APPROVED", "REJECTED"],
    default: "PENDING",
  })
  status!: string;
}
