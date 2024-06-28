import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
// import { Active, DataRef, Over } from "@dnd-kit/core";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function hasDraggableData() {}
