// Invoice Types - Re-exports from Wails-generated DTOs
// Source of truth: wailsjs/go/models.ts (auto-generated from Go dto package)

import { dto } from "@/wailsjs/go/models";

// Re-export Wails DTO types with frontend-friendly aliases
export type Invoice = dto.InvoiceOutput;
export type InvoiceItem = dto.InvoiceItemOutput;
export type InvoiceItemInput = dto.InvoiceItemInput;
export type CreateInvoiceInput = dto.CreateInvoiceInput;
export type UpdateInvoiceInput = dto.UpdateInvoiceInput;

// Service Interface
export interface IInvoiceService {
  list(): Promise<Invoice[]>;
  get(id: number): Promise<Invoice>;
  create(input: CreateInvoiceInput): Promise<Invoice>;
  update(input: UpdateInvoiceInput): Promise<Invoice>;
  delete(id: number): Promise<void>;
  generatePdf(id: number, message?: string): Promise<string>;
  sendEmail(id: number): Promise<boolean>;
  setTimeEntries(input: dto.SetInvoiceTimeEntriesInput): Promise<Invoice>;
}
