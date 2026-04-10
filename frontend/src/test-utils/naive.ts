import { vi } from "vitest";

export const mockMessage = {
  success: vi.fn(),
  error: vi.fn(),
  warning: vi.fn(),
  info: vi.fn(),
};

export function resetNaiveMocks() {
  mockMessage.success.mockReset();
  mockMessage.error.mockReset();
  mockMessage.warning.mockReset();
  mockMessage.info.mockReset();
}
