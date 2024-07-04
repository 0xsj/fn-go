/**
 * HOC for icons inside combobox
 */
import React from "react";
import type { SVGProps } from "react";

export interface SVGRProps {
  title?: string;
  titleId?: string;
  desc?: string;
  descId?: string;
}

type IconComponentProps = SVGProps<SVGSVGElement> & SVGRProps;

export const withIcon = (
  SvgContent: () => JSX.Element,
  fill?: string,
  ariaLabel?: string
) => {
  return ({ title, titleId, desc, descId, ...props }: IconComponentProps) => (
    <svg
      width='1em'
      height='1em'
      fill={fill ?? "#5e5e5f"}
      aria-label={ariaLabel}
      aria-labelledby={titleId}
      aria-describedby={descId}
      {...props}
    >
      {desc ? <desc id={descId}>{desc}</desc> : null}
      <title id={titleId}>{title}</title>
      <SvgContent />
    </svg>
  );
};
