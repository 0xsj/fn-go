import { cn } from "@/lib/utils";
import { HTMLAttributes } from "react";

interface Props extends HTMLAttributes<HTMLDivElement> {
  px?: number | string;
  py?: number | string;
  mx?: number | string;
  my?: number | string;
  m?: number | string;
  p?: number | string;
  center?: boolean;
  w?: string;
  h?: string;
  color?: string;
  bg?: string;
  maxW?: string;
  maxH?: string;
  pl?: number | string;
  pr?: number | string;
  pt?: number | string;
  pb?: number | string;
  ml?: number | string;
  mr?: number | string;
  mb?: number | string;
  mt?: number | string;
}
export const Box: React.FC<Props> = ({
  px,
  py,
  mx,
  my,
  m,
  p,
  center,
  w,
  h,
  color,
  bg,
  maxW,
  maxH,
  pl,
  pr,
  pt,
  pb,
  ml,
  mb,
  mt,
  mr,
  className,
  children,
  ...props
}) => {
  const classes = cn(
    `flex`,
    `px-${px}`,
    `py-${py}`,
    `mx-${mx}`,
    `my-${my}`,
    `m-${m}`,
    `p-${p}`,
    center && "text-center",
    `w-${w}`,
    `h-${h}`,
    `text-${color}`,
    `bg-${bg}`,
    `max-w-${maxW}`,
    `max-h-${maxH}`,
    `pl-${pl}`,
    `pr-${pr}`,
    `ml-${ml}`,
    `mr-${mr}`,
    `mb-${mb}`,
    className
  );

  return <div className={classes}>{children}</div>;
};
