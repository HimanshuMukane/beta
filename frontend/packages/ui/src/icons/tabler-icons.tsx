import React, { forwardRef } from "react";

export interface TablerIconProps extends Omit<React.SVGProps<SVGSVGElement>, "stroke"> {
  size?: number | string;
  stroke?: number | string;
  color?: string;
  className?: string;
}

export type Icon = React.ForwardRefExoticComponent<TablerIconProps & React.RefAttributes<SVGSVGElement>>;
export type TablerIcon = Icon;
export type IconProps = TablerIconProps;

export const IconAdjustments: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconAdjustments(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 10a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" />
  <path d="M6 4v4" />
  <path d="M6 12v8" />
  <path d="M10 16a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" />
  <path d="M12 4v10" />
  <path d="M12 18v2" />
  <path d="M16 7a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" />
  <path d="M18 4v1" />
  <path d="M18 9v11" />
    </svg>
  );
});
IconAdjustments.displayName = "IconAdjustments";

export const IconAdjustmentsHorizontal: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconAdjustmentsHorizontal(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 6a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" />
  <path d="M4 6l8 0" />
  <path d="M16 6l4 0" />
  <path d="M6 12a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" />
  <path d="M4 12l2 0" />
  <path d="M10 12l10 0" />
  <path d="M15 18a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" />
  <path d="M4 18l11 0" />
  <path d="M19 18l1 0" />
    </svg>
  );
});
IconAdjustmentsHorizontal.displayName = "IconAdjustmentsHorizontal";

export const IconAlertCircle: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconAlertCircle(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0" />
  <path d="M12 8v4" />
  <path d="M12 16h.01" />
    </svg>
  );
});
IconAlertCircle.displayName = "IconAlertCircle";

export const IconAlertTriangle: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconAlertTriangle(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 9v4" />
  <path d="M10.363 3.591l-8.106 13.534a1.914 1.914 0 0 0 1.636 2.871h16.214a1.914 1.914 0 0 0 1.636 -2.87l-8.106 -13.536a1.914 1.914 0 0 0 -3.274 0" />
  <path d="M12 16h.01" />
    </svg>
  );
});
IconAlertTriangle.displayName = "IconAlertTriangle";

export const IconAlignCenter: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconAlignCenter(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 6l16 0" />
  <path d="M8 12l8 0" />
  <path d="M6 18l12 0" />
    </svg>
  );
});
IconAlignCenter.displayName = "IconAlignCenter";

export const IconAlignLeft: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconAlignLeft(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 6l16 0" />
  <path d="M4 12l10 0" />
  <path d="M4 18l14 0" />
    </svg>
  );
});
IconAlignLeft.displayName = "IconAlignLeft";

export const IconAlignRight: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconAlignRight(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 6l16 0" />
  <path d="M10 12l10 0" />
  <path d="M6 18l14 0" />
    </svg>
  );
});
IconAlignRight.displayName = "IconAlignRight";

export const IconArchive: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconArchive(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 6a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2" />
  <path d="M5 8v10a2 2 0 0 0 2 2h10a2 2 0 0 0 2 -2v-10" />
  <path d="M10 12l4 0" />
    </svg>
  );
});
IconArchive.displayName = "IconArchive";

export const IconArrowBackUp: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconArrowBackUp(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M9 14l-4 -4l4 -4" />
  <path d="M5 10h11a4 4 0 1 1 0 8h-1" />
    </svg>
  );
});
IconArrowBackUp.displayName = "IconArrowBackUp";

export const IconArrowForwardUp: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconArrowForwardUp(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M15 14l4 -4l-4 -4" />
  <path d="M19 10h-11a4 4 0 1 0 0 8h1" />
    </svg>
  );
});
IconArrowForwardUp.displayName = "IconArrowForwardUp";

export const IconArrowLeft: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconArrowLeft(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 12l14 0" />
  <path d="M5 12l6 6" />
  <path d="M5 12l6 -6" />
    </svg>
  );
});
IconArrowLeft.displayName = "IconArrowLeft";

export const IconArrowRight: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconArrowRight(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 12l14 0" />
  <path d="M13 18l6 -6" />
  <path d="M13 6l6 6" />
    </svg>
  );
});
IconArrowRight.displayName = "IconArrowRight";

export const IconArrowUpRight: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconArrowUpRight(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M17 7l-10 10" />
  <path d="M8 7l9 0l0 9" />
    </svg>
  );
});
IconArrowUpRight.displayName = "IconArrowUpRight";

export const IconArrowsSort: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconArrowsSort(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 9l4 -4l4 4m-4 -4v14" />
  <path d="M21 15l-4 4l-4 -4m4 4v-14" />
    </svg>
  );
});
IconArrowsSort.displayName = "IconArrowsSort";

export const IconArrowsUpDown: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconArrowsUpDown(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M7 3l0 18" />
  <path d="M10 6l-3 -3l-3 3" />
  <path d="M20 18l-3 3l-3 -3" />
  <path d="M17 21l0 -18" />
    </svg>
  );
});
IconArrowsUpDown.displayName = "IconArrowsUpDown";

export const IconBell: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconBell(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M10 5a2 2 0 1 1 4 0a7 7 0 0 1 4 6v3a4 4 0 0 0 2 3h-16a4 4 0 0 0 2 -3v-3a7 7 0 0 1 4 -6" />
  <path d="M9 17v1a3 3 0 0 0 6 0v-1" />
    </svg>
  );
});
IconBell.displayName = "IconBell";

export const IconBlockquote: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconBlockquote(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M6 15h15" />
  <path d="M21 19h-15" />
  <path d="M15 11h6" />
  <path d="M21 7h-6" />
  <path d="M9 9h1a1 1 0 1 1 -1 1v-2.5a2 2 0 0 1 2 -2" />
  <path d="M3 9h1a1 1 0 1 1 -1 1v-2.5a2 2 0 0 1 2 -2" />
    </svg>
  );
});
IconBlockquote.displayName = "IconBlockquote";

export const IconBold: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconBold(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M7 5h6a3.5 3.5 0 0 1 0 7h-6l0 -7" />
  <path d="M13 12h1a3.5 3.5 0 0 1 0 7h-7v-7" />
    </svg>
  );
});
IconBold.displayName = "IconBold";

export const IconBook: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconBook(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 19a9 9 0 0 1 9 0a9 9 0 0 1 9 0" />
  <path d="M3 6a9 9 0 0 1 9 0a9 9 0 0 1 9 0" />
  <path d="M3 6l0 13" />
  <path d="M12 6l0 13" />
  <path d="M21 6l0 13" />
    </svg>
  );
});
IconBook.displayName = "IconBook";

export const IconBox: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconBox(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 3l8 4.5l0 9l-8 4.5l-8 -4.5l0 -9l8 -4.5" />
  <path d="M12 12l8 -4.5" />
  <path d="M12 12l0 9" />
  <path d="M12 12l-8 -4.5" />
    </svg>
  );
});
IconBox.displayName = "IconBox";

export const IconBriefcase: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconBriefcase(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 9a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v9a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2l0 -9" />
  <path d="M8 7v-2a2 2 0 0 1 2 -2h4a2 2 0 0 1 2 2v2" />
  <path d="M12 12l0 .01" />
  <path d="M3 13a20 20 0 0 0 18 0" />
    </svg>
  );
});
IconBriefcase.displayName = "IconBriefcase";

export const IconCalendar: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCalendar(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 7a2 2 0 0 1 2 -2h12a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2v-12" />
  <path d="M16 3v4" />
  <path d="M8 3v4" />
  <path d="M4 11h16" />
    </svg>
  );
});
IconCalendar.displayName = "IconCalendar";

export const IconCalendarEvent: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCalendarEvent(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 7a2 2 0 0 1 2 -2h12a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2l0 -12" />
  <path d="M16 3l0 4" />
  <path d="M8 3l0 4" />
  <path d="M4 11l16 0" />
  <path d="M8 15h2v2h-2l0 -2" />
    </svg>
  );
});
IconCalendarEvent.displayName = "IconCalendarEvent";

export const IconCamera: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCamera(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 7h1a2 2 0 0 0 2 -2a1 1 0 0 1 1 -1h6a1 1 0 0 1 1 1a2 2 0 0 0 2 2h1a2 2 0 0 1 2 2v9a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2v-9a2 2 0 0 1 2 -2" />
  <path d="M9 13a3 3 0 1 0 6 0a3 3 0 0 0 -6 0" />
    </svg>
  );
});
IconCamera.displayName = "IconCamera";

export const IconCash: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCash(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M7 15h-3a1 1 0 0 1 -1 -1v-8a1 1 0 0 1 1 -1h12a1 1 0 0 1 1 1v3" />
  <path d="M7 10a1 1 0 0 1 1 -1h12a1 1 0 0 1 1 1v8a1 1 0 0 1 -1 1h-12a1 1 0 0 1 -1 -1l0 -8" />
  <path d="M12 14a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" />
    </svg>
  );
});
IconCash.displayName = "IconCash";

export const IconChartBar: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconChartBar(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 13a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v6a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1l0 -6" />
  <path d="M15 9a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v10a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1l0 -10" />
  <path d="M9 5a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v14a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1l0 -14" />
  <path d="M4 20h14" />
    </svg>
  );
});
IconChartBar.displayName = "IconChartBar";

export const IconChartPie: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconChartPie(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M10 3.2a9 9 0 1 0 10.8 10.8a1 1 0 0 0 -1 -1h-6.8a2 2 0 0 1 -2 -2v-7a.9 .9 0 0 0 -1 -.8" />
  <path d="M15 3.5a9 9 0 0 1 5.5 5.5h-4.5a1 1 0 0 1 -1 -1v-4.5" />
    </svg>
  );
});
IconChartPie.displayName = "IconChartPie";

export const IconCheck: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCheck(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 12l5 5l10 -10" />
    </svg>
  );
});
IconCheck.displayName = "IconCheck";

export const IconChecklist: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconChecklist(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M9.615 20h-2.615a2 2 0 0 1 -2 -2v-12a2 2 0 0 1 2 -2h8a2 2 0 0 1 2 2v8" />
  <path d="M14 19l2 2l4 -4" />
  <path d="M9 8h4" />
  <path d="M9 12h2" />
    </svg>
  );
});
IconChecklist.displayName = "IconChecklist";

export const IconChevronDown: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconChevronDown(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M6 9l6 6l6 -6" />
    </svg>
  );
});
IconChevronDown.displayName = "IconChevronDown";

export const IconChevronLeft: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconChevronLeft(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M15 6l-6 6l6 6" />
    </svg>
  );
});
IconChevronLeft.displayName = "IconChevronLeft";

export const IconChevronRight: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconChevronRight(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M9 6l6 6l-6 6" />
    </svg>
  );
});
IconChevronRight.displayName = "IconChevronRight";

export const IconChevronUp: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconChevronUp(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M6 15l6 -6l6 6" />
    </svg>
  );
});
IconChevronUp.displayName = "IconChevronUp";

export const IconChevronsLeft: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconChevronsLeft(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M11 7l-5 5l5 5" />
  <path d="M17 7l-5 5l5 5" />
    </svg>
  );
});
IconChevronsLeft.displayName = "IconChevronsLeft";

export const IconChevronsRight: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconChevronsRight(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M7 7l5 5l-5 5" />
  <path d="M13 7l5 5l-5 5" />
    </svg>
  );
});
IconChevronsRight.displayName = "IconChevronsRight";

export const IconCircleCheck: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCircleCheck(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 12a9 9 0 1 0 18 0a9 9 0 1 0 -18 0" />
  <path d="M9 12l2 2l4 -4" />
    </svg>
  );
});
IconCircleCheck.displayName = "IconCircleCheck";

export const IconCircleDot: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCircleDot(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M11 12a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M3 12a9 9 0 1 0 18 0a9 9 0 1 0 -18 0" />
    </svg>
  );
});
IconCircleDot.displayName = "IconCircleDot";

export const IconCircleFilled: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCircleFilled(
  { size = 24, stroke = 0, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="currentColor"
      stroke={color || "none"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M7 3.34a10 10 0 1 1 -4.995 8.984l-.005 -.324l.005 -.324a10 10 0 0 1 4.995 -8.336z" />
    </svg>
  );
});
IconCircleFilled.displayName = "IconCircleFilled";

export const IconCircleX: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCircleX(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 12a9 9 0 1 0 18 0a9 9 0 1 0 -18 0" />
  <path d="M10 10l4 4m0 -4l-4 4" />
    </svg>
  );
});
IconCircleX.displayName = "IconCircleX";

export const IconClock: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconClock(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0" />
  <path d="M12 7v5l3 3" />
    </svg>
  );
});
IconClock.displayName = "IconClock";

export const IconClockHour4: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconClockHour4(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 12a9 9 0 1 0 18 0a9 9 0 1 0 -18 0" />
  <path d="M12 12l3 2" />
  <path d="M12 7v5" />
    </svg>
  );
});
IconClockHour4.displayName = "IconClockHour4";

export const IconCode: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCode(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M7 8l-4 4l4 4" />
  <path d="M17 8l4 4l-4 4" />
  <path d="M14 4l-4 16" />
    </svg>
  );
});
IconCode.displayName = "IconCode";

export const IconCoins: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCoins(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M9 14c0 1.657 2.686 3 6 3s6 -1.343 6 -3s-2.686 -3 -6 -3s-6 1.343 -6 3" />
  <path d="M9 14v4c0 1.656 2.686 3 6 3s6 -1.344 6 -3v-4" />
  <path d="M3 6c0 1.072 1.144 2.062 3 2.598s4.144 .536 6 0c1.856 -.536 3 -1.526 3 -2.598c0 -1.072 -1.144 -2.062 -3 -2.598s-4.144 -.536 -6 0c-1.856 .536 -3 1.526 -3 2.598" />
  <path d="M3 6v10c0 .888 .772 1.45 2 2" />
  <path d="M3 11c0 .888 .772 1.45 2 2" />
    </svg>
  );
});
IconCoins.displayName = "IconCoins";

export const IconCommand: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCommand(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M7 9a2 2 0 1 1 2 -2v10a2 2 0 1 1 -2 -2h10a2 2 0 1 1 -2 2v-10a2 2 0 1 1 2 2h-10" />
    </svg>
  );
});
IconCommand.displayName = "IconCommand";

export const IconCopy: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCopy(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M7 9.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667l0 -8.666" />
  <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1" />
    </svg>
  );
});
IconCopy.displayName = "IconCopy";

export const IconCornerDownRight: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCornerDownRight(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M6 6v6a3 3 0 0 0 3 3h10l-4 -4m0 8l4 -4" />
    </svg>
  );
});
IconCornerDownRight.displayName = "IconCornerDownRight";

export const IconCreditCard: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCreditCard(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 8a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v8a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3l0 -8" />
  <path d="M3 10l18 0" />
  <path d="M7 15l.01 0" />
  <path d="M11 15l2 0" />
    </svg>
  );
});
IconCreditCard.displayName = "IconCreditCard";

export const IconCurrencyRupee: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconCurrencyRupee(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M18 5h-11h3a4 4 0 0 1 0 8h-3l6 6" />
  <path d="M7 9l11 0" />
    </svg>
  );
});
IconCurrencyRupee.displayName = "IconCurrencyRupee";

export const IconDatabase: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconDatabase(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 6a8 3 0 1 0 16 0a8 3 0 1 0 -16 0" />
  <path d="M4 6v6a8 3 0 0 0 16 0v-6" />
  <path d="M4 12v6a8 3 0 0 0 16 0v-6" />
    </svg>
  );
});
IconDatabase.displayName = "IconDatabase";

export const IconDeviceFloppy: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconDeviceFloppy(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M6 4h10l4 4v10a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2v-12a2 2 0 0 1 2 -2" />
  <path d="M10 14a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" />
  <path d="M14 4l0 4l-6 0l0 -4" />
    </svg>
  );
});
IconDeviceFloppy.displayName = "IconDeviceFloppy";

export const IconDots: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconDots(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 12a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M11 12a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M18 12a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
    </svg>
  );
});
IconDots.displayName = "IconDots";

export const IconDotsVertical: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconDotsVertical(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M11 12a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M11 19a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M11 5a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
    </svg>
  );
});
IconDotsVertical.displayName = "IconDotsVertical";

export const IconDownload: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconDownload(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 17v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2 -2v-2" />
  <path d="M7 11l5 5l5 -5" />
  <path d="M12 4l0 12" />
    </svg>
  );
});
IconDownload.displayName = "IconDownload";

export const IconEdit: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconEdit(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M7 7h-1a2 2 0 0 0 -2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2 -2v-1" />
  <path d="M20.385 6.585a2.1 2.1 0 0 0 -2.97 -2.97l-8.415 8.385v3h3l8.385 -8.415" />
  <path d="M16 5l3 3" />
    </svg>
  );
});
IconEdit.displayName = "IconEdit";

export const IconExternalLink: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconExternalLink(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 6h-6a2 2 0 0 0 -2 2v10a2 2 0 0 0 2 2h10a2 2 0 0 0 2 -2v-6" />
  <path d="M11 13l9 -9" />
  <path d="M15 4h5v5" />
    </svg>
  );
});
IconExternalLink.displayName = "IconExternalLink";

export const IconEye: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconEye(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" />
  <path d="M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6" />
    </svg>
  );
});
IconEye.displayName = "IconEye";

export const IconFileCheck: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconFileCheck(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M14 3v4a1 1 0 0 0 1 1h4" />
  <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2" />
  <path d="M9 15l2 2l4 -4" />
    </svg>
  );
});
IconFileCheck.displayName = "IconFileCheck";

export const IconFileInvoice: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconFileInvoice(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M14 3v4a1 1 0 0 0 1 1h4" />
  <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2" />
  <path d="M9 7l1 0" />
  <path d="M9 13l6 0" />
  <path d="M13 17l2 0" />
    </svg>
  );
});
IconFileInvoice.displayName = "IconFileInvoice";

export const IconFileSpreadsheet: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconFileSpreadsheet(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M14 3v4a1 1 0 0 0 1 1h4" />
  <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2" />
  <path d="M8 11h8v7h-8l0 -7" />
  <path d="M8 15h8" />
  <path d="M11 11v7" />
    </svg>
  );
});
IconFileSpreadsheet.displayName = "IconFileSpreadsheet";

export const IconFileText: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconFileText(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M14 3v4a1 1 0 0 0 1 1h4" />
  <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2" />
  <path d="M9 9l1 0" />
  <path d="M9 13l6 0" />
  <path d="M9 17l6 0" />
    </svg>
  );
});
IconFileText.displayName = "IconFileText";

export const IconFilter: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconFilter(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 4h16v2.172a2 2 0 0 1 -.586 1.414l-4.414 4.414v7l-6 2v-8.5l-4.48 -4.928a2 2 0 0 1 -.52 -1.345v-2.227" />
    </svg>
  );
});
IconFilter.displayName = "IconFilter";

export const IconFlame: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconFlame(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 10.941c2.333 -3.308 .167 -7.823 -1 -8.941c0 3.395 -2.235 5.299 -3.667 6.706c-1.43 1.408 -2.333 3.294 -2.333 5.588c0 3.704 3.134 6.706 7 6.706c3.866 0 7 -3.002 7 -6.706c0 -1.712 -1.232 -4.403 -2.333 -5.588c-2.084 3.353 -3.257 3.353 -4.667 2.235" />
    </svg>
  );
});
IconFlame.displayName = "IconFlame";

export const IconFolder: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconFolder(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 4h4l3 3h7a2 2 0 0 1 2 2v8a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2v-11a2 2 0 0 1 2 -2" />
    </svg>
  );
});
IconFolder.displayName = "IconFolder";

export const IconFolderPlus: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconFolderPlus(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 19h-7a2 2 0 0 1 -2 -2v-11a2 2 0 0 1 2 -2h4l3 3h7a2 2 0 0 1 2 2v3.5" />
  <path d="M16 19h6" />
  <path d="M19 16v6" />
    </svg>
  );
});
IconFolderPlus.displayName = "IconFolderPlus";

export const IconFolders: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconFolders(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M9 3h3l2 2h5a2 2 0 0 1 2 2v7a2 2 0 0 1 -2 2h-10a2 2 0 0 1 -2 -2v-9a2 2 0 0 1 2 -2" />
  <path d="M17 16v2a2 2 0 0 1 -2 2h-10a2 2 0 0 1 -2 -2v-9a2 2 0 0 1 2 -2h2" />
    </svg>
  );
});
IconFolders.displayName = "IconFolders";

export const IconForms: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconForms(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 3a3 3 0 0 0 -3 3v12a3 3 0 0 0 3 3" />
  <path d="M6 3a3 3 0 0 1 3 3v12a3 3 0 0 1 -3 3" />
  <path d="M13 7h7a1 1 0 0 1 1 1v8a1 1 0 0 1 -1 1h-7" />
  <path d="M5 7h-1a1 1 0 0 0 -1 1v8a1 1 0 0 0 1 1h1" />
  <path d="M17 12h.01" />
  <path d="M13 12h.01" />
    </svg>
  );
});
IconForms.displayName = "IconForms";

export const IconGitFork: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconGitFork(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M10 18a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" />
  <path d="M5 6a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" />
  <path d="M15 6a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" />
  <path d="M7 8v2a2 2 0 0 0 2 2h6a2 2 0 0 0 2 -2v-2" />
  <path d="M12 12l0 4" />
    </svg>
  );
});
IconGitFork.displayName = "IconGitFork";

export const IconGripVertical: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconGripVertical(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M8 5a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M8 12a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M8 19a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M14 5a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M14 12a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M14 19a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
    </svg>
  );
});
IconGripVertical.displayName = "IconGripVertical";

export const IconH1: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconH1(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M19 18v-8l-2 2" />
  <path d="M4 6v12" />
  <path d="M12 6v12" />
  <path d="M11 18h2" />
  <path d="M3 18h2" />
  <path d="M4 12h8" />
  <path d="M3 6h2" />
  <path d="M11 6h2" />
    </svg>
  );
});
IconH1.displayName = "IconH1";

export const IconH2: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconH2(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M17 12a2 2 0 1 1 4 0c0 .591 -.417 1.318 -.816 1.858l-3.184 4.143l4 0" />
  <path d="M4 6v12" />
  <path d="M12 6v12" />
  <path d="M11 18h2" />
  <path d="M3 18h2" />
  <path d="M4 12h8" />
  <path d="M3 6h2" />
  <path d="M11 6h2" />
    </svg>
  );
});
IconH2.displayName = "IconH2";

export const IconHome: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconHome(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 12l-2 0l9 -9l9 9l-2 0" />
  <path d="M5 12v7a2 2 0 0 0 2 2h10a2 2 0 0 0 2 -2v-7" />
  <path d="M9 21v-6a2 2 0 0 1 2 -2h2a2 2 0 0 1 2 2v6" />
    </svg>
  );
});
IconHome.displayName = "IconHome";

export const IconId: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconId(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 7a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v10a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3l0 -10" />
  <path d="M7 10a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" />
  <path d="M15 8l2 0" />
  <path d="M15 12l2 0" />
  <path d="M7 16l10 0" />
    </svg>
  );
});
IconId.displayName = "IconId";

export const IconInbox: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconInbox(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 6a2 2 0 0 1 2 -2h12a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2l0 -12" />
  <path d="M4 13h3l3 3h4l3 -3h3" />
    </svg>
  );
});
IconInbox.displayName = "IconInbox";

export const IconInfoCircle: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconInfoCircle(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0" />
  <path d="M12 9h.01" />
  <path d="M11 12h1v4h1" />
    </svg>
  );
});
IconInfoCircle.displayName = "IconInfoCircle";

export const IconItalic: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconItalic(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M11 5l6 0" />
  <path d="M7 19l6 0" />
  <path d="M14 5l-4 14" />
    </svg>
  );
});
IconItalic.displayName = "IconItalic";

export const IconLayersSubtract: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconLayersSubtract(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M8 6a2 2 0 0 1 2 -2h8a2 2 0 0 1 2 2v8a2 2 0 0 1 -2 2h-8a2 2 0 0 1 -2 -2l0 -8" />
  <path d="M16 16v2a2 2 0 0 1 -2 2h-8a2 2 0 0 1 -2 -2v-8a2 2 0 0 1 2 -2h2" />
    </svg>
  );
});
IconLayersSubtract.displayName = "IconLayersSubtract";

export const IconLayout: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconLayout(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 6a2 2 0 0 1 2 -2h2a2 2 0 0 1 2 2v1a2 2 0 0 1 -2 2h-2a2 2 0 0 1 -2 -2l0 -1" />
  <path d="M4 15a2 2 0 0 1 2 -2h2a2 2 0 0 1 2 2v3a2 2 0 0 1 -2 2h-2a2 2 0 0 1 -2 -2l0 -3" />
  <path d="M14 6a2 2 0 0 1 2 -2h2a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-2a2 2 0 0 1 -2 -2l0 -12" />
    </svg>
  );
});
IconLayout.displayName = "IconLayout";

export const IconLayoutKanban: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconLayoutKanban(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 4l6 0" />
  <path d="M14 4l6 0" />
  <path d="M4 10a2 2 0 0 1 2 -2h2a2 2 0 0 1 2 2v8a2 2 0 0 1 -2 2h-2a2 2 0 0 1 -2 -2l0 -8" />
  <path d="M14 10a2 2 0 0 1 2 -2h2a2 2 0 0 1 2 2v2a2 2 0 0 1 -2 2h-2a2 2 0 0 1 -2 -2l0 -2" />
    </svg>
  );
});
IconLayoutKanban.displayName = "IconLayoutKanban";

export const IconLayoutSidebar: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconLayoutSidebar(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 6a2 2 0 0 1 2 -2h12a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2l0 -12" />
  <path d="M9 4l0 16" />
    </svg>
  );
});
IconLayoutSidebar.displayName = "IconLayoutSidebar";

export const IconList: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconList(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M9 6l11 0" />
  <path d="M9 12l11 0" />
  <path d="M9 18l11 0" />
  <path d="M5 6l0 .01" />
  <path d="M5 12l0 .01" />
  <path d="M5 18l0 .01" />
    </svg>
  );
});
IconList.displayName = "IconList";

export const IconListCheck: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconListCheck(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3.5 5.5l1.5 1.5l2.5 -2.5" />
  <path d="M3.5 11.5l1.5 1.5l2.5 -2.5" />
  <path d="M3.5 17.5l1.5 1.5l2.5 -2.5" />
  <path d="M11 6l9 0" />
  <path d="M11 12l9 0" />
  <path d="M11 18l9 0" />
    </svg>
  );
});
IconListCheck.displayName = "IconListCheck";

export const IconListNumbers: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconListNumbers(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M11 6h9" />
  <path d="M11 12h9" />
  <path d="M12 18h8" />
  <path d="M4 16a2 2 0 1 1 4 0c0 .591 -.5 1 -1 1.5l-3 2.5h4" />
  <path d="M6 10v-6l-2 2" />
    </svg>
  );
});
IconListNumbers.displayName = "IconListNumbers";

export const IconLoader2: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconLoader2(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 3a9 9 0 1 0 9 9" />
    </svg>
  );
});
IconLoader2.displayName = "IconLoader2";

export const IconLock: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconLock(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 13a2 2 0 0 1 2 -2h10a2 2 0 0 1 2 2v6a2 2 0 0 1 -2 2h-10a2 2 0 0 1 -2 -2v-6" />
  <path d="M11 16a1 1 0 1 0 2 0a1 1 0 0 0 -2 0" />
  <path d="M8 11v-4a4 4 0 1 1 8 0v4" />
    </svg>
  );
});
IconLock.displayName = "IconLock";

export const IconLogout: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconLogout(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M14 8v-2a2 2 0 0 0 -2 -2h-7a2 2 0 0 0 -2 2v12a2 2 0 0 0 2 2h7a2 2 0 0 0 2 -2v-2" />
  <path d="M9 12h12l-3 -3" />
  <path d="M18 15l3 -3" />
    </svg>
  );
});
IconLogout.displayName = "IconLogout";

export const IconMail: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconMail(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 7a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v10a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2v-10" />
  <path d="M3 7l9 6l9 -6" />
    </svg>
  );
});
IconMail.displayName = "IconMail";

export const IconMap: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconMap(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 7l6 -3l6 3l6 -3v13l-6 3l-6 -3l-6 3v-13" />
  <path d="M9 4v13" />
  <path d="M15 7v13" />
    </svg>
  );
});
IconMap.displayName = "IconMap";

export const IconMapPin: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconMapPin(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M9 11a3 3 0 1 0 6 0a3 3 0 0 0 -6 0" />
  <path d="M17.657 16.657l-4.243 4.243a2 2 0 0 1 -2.827 0l-4.244 -4.243a8 8 0 1 1 11.314 0" />
    </svg>
  );
});
IconMapPin.displayName = "IconMapPin";

export const IconMessageCircle: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconMessageCircle(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 20l1.3 -3.9c-2.324 -3.437 -1.426 -7.872 2.1 -10.374c3.526 -2.501 8.59 -2.296 11.845 .48c3.255 2.777 3.695 7.266 1.029 10.501c-2.666 3.235 -7.615 4.215 -11.574 2.293l-4.7 1" />
    </svg>
  );
});
IconMessageCircle.displayName = "IconMessageCircle";

export const IconMessageDots: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconMessageDots(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 11v.01" />
  <path d="M8 11v.01" />
  <path d="M16 11v.01" />
  <path d="M18 4a3 3 0 0 1 3 3v8a3 3 0 0 1 -3 3h-5l-5 3v-3h-2a3 3 0 0 1 -3 -3v-8a3 3 0 0 1 3 -3l12 0" />
    </svg>
  );
});
IconMessageDots.displayName = "IconMessageDots";

export const IconMessageHeart: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconMessageHeart(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M8 9h8" />
  <path d="M8 13h3.5" />
  <path d="M10.48 19.512l-2.48 1.488v-3h-2a3 3 0 0 1 -3 -3v-8a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v4" />
  <path d="M18 22l3.35 -3.284a2.143 2.143 0 0 0 .005 -3.071a2.242 2.242 0 0 0 -3.129 -.006l-.224 .22l-.223 -.22a2.242 2.242 0 0 0 -3.128 -.006a2.143 2.143 0 0 0 -.006 3.071l3.355 3.296" />
    </svg>
  );
});
IconMessageHeart.displayName = "IconMessageHeart";

export const IconMinus: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconMinus(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 12l14 0" />
    </svg>
  );
});
IconMinus.displayName = "IconMinus";

export const IconMoodSmile: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconMoodSmile(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 12a9 9 0 1 0 18 0a9 9 0 1 0 -18 0" />
  <path d="M9 10l.01 0" />
  <path d="M15 10l.01 0" />
  <path d="M9.5 15a3.5 3.5 0 0 0 5 0" />
    </svg>
  );
});
IconMoodSmile.displayName = "IconMoodSmile";

export const IconMoon: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconMoon(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 3c.132 0 .263 0 .393 0a7.5 7.5 0 0 0 7.92 12.446a9 9 0 1 1 -8.313 -12.454l0 .008" />
    </svg>
  );
});
IconMoon.displayName = "IconMoon";

export const IconNote: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconNote(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M13 20l7 -7" />
  <path d="M13 20v-6a1 1 0 0 1 1 -1h6v-7a2 2 0 0 0 -2 -2h-12a2 2 0 0 0 -2 2v12a2 2 0 0 0 2 2h7" />
    </svg>
  );
});
IconNote.displayName = "IconNote";

export const IconNotes: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconNotes(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 5a2 2 0 0 1 2 -2h10a2 2 0 0 1 2 2v14a2 2 0 0 1 -2 2h-10a2 2 0 0 1 -2 -2l0 -14" />
  <path d="M9 7l6 0" />
  <path d="M9 11l6 0" />
  <path d="M9 15l4 0" />
    </svg>
  );
});
IconNotes.displayName = "IconNotes";

export const IconPackage: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconPackage(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 3l8 4.5l0 9l-8 4.5l-8 -4.5l0 -9l8 -4.5" />
  <path d="M12 12l8 -4.5" />
  <path d="M12 12l0 9" />
  <path d="M12 12l-8 -4.5" />
  <path d="M16 5.25l-8 4.5" />
    </svg>
  );
});
IconPackage.displayName = "IconPackage";

export const IconPalette: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconPalette(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 21a9 9 0 0 1 0 -18c4.97 0 9 3.582 9 8c0 1.06 -.474 2.078 -1.318 2.828c-.844 .75 -1.989 1.172 -3.182 1.172h-2.5a2 2 0 0 0 -1 3.75a1.3 1.3 0 0 1 -1 2.25" />
  <path d="M7.5 10.5a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M11.5 7.5a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M15.5 10.5a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
    </svg>
  );
});
IconPalette.displayName = "IconPalette";

export const IconPaperclip: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconPaperclip(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M15 7l-6.5 6.5a1.5 1.5 0 0 0 3 3l6.5 -6.5a3 3 0 0 0 -6 -6l-6.5 6.5a4.5 4.5 0 0 0 9 9l6.5 -6.5" />
    </svg>
  );
});
IconPaperclip.displayName = "IconPaperclip";

export const IconPercentage: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconPercentage(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M16 17a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M6 7a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M6 18l12 -12" />
    </svg>
  );
});
IconPercentage.displayName = "IconPercentage";

export const IconPhone: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconPhone(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 4h4l2 5l-2.5 1.5a11 11 0 0 0 5 5l1.5 -2.5l5 2v4a2 2 0 0 1 -2 2a16 16 0 0 1 -15 -15a2 2 0 0 1 2 -2" />
    </svg>
  );
});
IconPhone.displayName = "IconPhone";

export const IconPhoto: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconPhoto(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M15 8h.01" />
  <path d="M3 6a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v12a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3v-12" />
  <path d="M3 16l5 -5c.928 -.893 2.072 -.893 3 0l5 5" />
  <path d="M14 14l1 -1c.928 -.893 2.072 -.893 3 0l3 3" />
    </svg>
  );
});
IconPhoto.displayName = "IconPhoto";

export const IconPlus: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconPlus(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 5l0 14" />
  <path d="M5 12l14 0" />
    </svg>
  );
});
IconPlus.displayName = "IconPlus";

export const IconPrinter: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconPrinter(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M17 17h2a2 2 0 0 0 2 -2v-4a2 2 0 0 0 -2 -2h-14a2 2 0 0 0 -2 2v4a2 2 0 0 0 2 2h2" />
  <path d="M17 9v-4a2 2 0 0 0 -2 -2h-6a2 2 0 0 0 -2 2v4" />
  <path d="M7 15a2 2 0 0 1 2 -2h6a2 2 0 0 1 2 2v4a2 2 0 0 1 -2 2h-6a2 2 0 0 1 -2 -2l0 -4" />
    </svg>
  );
});
IconPrinter.displayName = "IconPrinter";

export const IconReceipt: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconReceipt(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 21v-16a2 2 0 0 1 2 -2h10a2 2 0 0 1 2 2v16l-3 -2l-2 2l-2 -2l-2 2l-2 -2l-3 2m4 -14h6m-6 4h6m-2 4h2" />
    </svg>
  );
});
IconReceipt.displayName = "IconReceipt";

export const IconReceipt2: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconReceipt2(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 21v-16a2 2 0 0 1 2 -2h10a2 2 0 0 1 2 2v16l-3 -2l-2 2l-2 -2l-2 2l-2 -2l-3 2" />
  <path d="M14 8h-2.5a1.5 1.5 0 0 0 0 3h1a1.5 1.5 0 0 1 0 3h-2.5m2 0v1.5m0 -9v1.5" />
    </svg>
  );
});
IconReceipt2.displayName = "IconReceipt2";

export const IconRefresh: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconRefresh(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M20 11a8.1 8.1 0 0 0 -15.5 -2m-.5 -4v4h4" />
  <path d="M4 13a8.1 8.1 0 0 0 15.5 2m.5 4v-4h-4" />
    </svg>
  );
});
IconRefresh.displayName = "IconRefresh";

export const IconReportAnalytics: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconReportAnalytics(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M9 5h-2a2 2 0 0 0 -2 2v12a2 2 0 0 0 2 2h10a2 2 0 0 0 2 -2v-12a2 2 0 0 0 -2 -2h-2" />
  <path d="M9 5a2 2 0 0 1 2 -2h2a2 2 0 0 1 2 2a2 2 0 0 1 -2 2h-2a2 2 0 0 1 -2 -2" />
  <path d="M9 17v-5" />
  <path d="M12 17v-1" />
  <path d="M15 17v-3" />
    </svg>
  );
});
IconReportAnalytics.displayName = "IconReportAnalytics";

export const IconRobot: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconRobot(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M6 6a2 2 0 0 1 2 -2h8a2 2 0 0 1 2 2v4a2 2 0 0 1 -2 2h-8a2 2 0 0 1 -2 -2l0 -4" />
  <path d="M12 2v2" />
  <path d="M9 12v9" />
  <path d="M15 12v9" />
  <path d="M5 16l4 -2" />
  <path d="M15 14l4 2" />
  <path d="M9 18h6" />
  <path d="M10 8v.01" />
  <path d="M14 8v.01" />
    </svg>
  );
});
IconRobot.displayName = "IconRobot";

export const IconRotate: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconRotate(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M19.95 11a8 8 0 1 0 -.5 4m.5 5v-5h-5" />
    </svg>
  );
});
IconRotate.displayName = "IconRotate";

export const IconRotateClockwise2: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconRotateClockwise2(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M9 4.55a8 8 0 0 1 6 14.9m0 -4.45v5h5" />
  <path d="M5.63 7.16l0 .01" />
  <path d="M4.06 11l0 .01" />
  <path d="M4.63 15.1l0 .01" />
  <path d="M7.16 18.37l0 .01" />
  <path d="M11 19.94l0 .01" />
    </svg>
  );
});
IconRotateClockwise2.displayName = "IconRotateClockwise2";

export const IconSearch: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconSearch(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 10a7 7 0 1 0 14 0a7 7 0 1 0 -14 0" />
  <path d="M21 21l-6 -6" />
    </svg>
  );
});
IconSearch.displayName = "IconSearch";

export const IconSelector: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconSelector(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M8 9l4 -4l4 4" />
  <path d="M16 15l-4 4l-4 -4" />
    </svg>
  );
});
IconSelector.displayName = "IconSelector";

export const IconSeparator: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconSeparator(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 12l0 .01" />
  <path d="M7 12l10 0" />
  <path d="M21 12l0 .01" />
    </svg>
  );
});
IconSeparator.displayName = "IconSeparator";

export const IconSettings: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconSettings(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065" />
  <path d="M9 12a3 3 0 1 0 6 0a3 3 0 0 0 -6 0" />
    </svg>
  );
});
IconSettings.displayName = "IconSettings";

export const IconShare3: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconShare3(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M13 4v4c-6.575 1.028 -9.02 6.788 -10 12c-.037 .206 5.384 -5.962 10 -6v4l8 -7l-8 -7" />
    </svg>
  );
});
IconShare3.displayName = "IconShare3";

export const IconShoppingBag: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconShoppingBag(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M6.331 8h11.339a2 2 0 0 1 1.977 2.304l-1.255 8.152a3 3 0 0 1 -2.966 2.544h-6.852a3 3 0 0 1 -2.965 -2.544l-1.255 -8.152a2 2 0 0 1 1.977 -2.304" />
  <path d="M9 11v-5a3 3 0 0 1 6 0v5" />
    </svg>
  );
});
IconShoppingBag.displayName = "IconShoppingBag";

export const IconSignature: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconSignature(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 17c3.333 -3.333 5 -6 5 -8c0 -3 -1 -3 -2 -3s-2.032 1.085 -2 3c.034 2.048 1.658 4.877 2.5 6c1.5 2 2.5 2.5 3.5 1l2 -3c.333 2.667 1.333 4 3 4c.53 0 2.639 -2 3 -2c.517 0 1.517 .667 3 2" />
    </svg>
  );
});
IconSignature.displayName = "IconSignature";

export const IconSortAscending: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconSortAscending(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 6l7 0" />
  <path d="M4 12l7 0" />
  <path d="M4 18l9 0" />
  <path d="M15 9l3 -3l3 3" />
  <path d="M18 6l0 12" />
    </svg>
  );
});
IconSortAscending.displayName = "IconSortAscending";

export const IconSortDescending: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconSortDescending(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 6l9 0" />
  <path d="M4 12l7 0" />
  <path d="M4 18l7 0" />
  <path d="M15 15l3 3l3 -3" />
  <path d="M18 6l0 12" />
    </svg>
  );
});
IconSortDescending.displayName = "IconSortDescending";

export const IconSparkles: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconSparkles(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M16 18a2 2 0 0 1 2 2a2 2 0 0 1 2 -2a2 2 0 0 1 -2 -2a2 2 0 0 1 -2 2m0 -12a2 2 0 0 1 2 2a2 2 0 0 1 2 -2a2 2 0 0 1 -2 -2a2 2 0 0 1 -2 2m-7 12a6 6 0 0 1 6 -6a6 6 0 0 1 -6 -6a6 6 0 0 1 -6 6a6 6 0 0 1 6 6" />
    </svg>
  );
});
IconSparkles.displayName = "IconSparkles";

export const IconSquarePlus: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconSquarePlus(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M9 12h6" />
  <path d="M12 9v6" />
  <path d="M3 5a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v14a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2v-14" />
    </svg>
  );
});
IconSquarePlus.displayName = "IconSquarePlus";

export const IconStack2: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconStack2(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M12 4l-8 4l8 4l8 -4l-8 -4" />
  <path d="M4 12l8 4l8 -4" />
  <path d="M4 16l8 4l8 -4" />
    </svg>
  );
});
IconStack2.displayName = "IconStack2";

export const IconStarFilled: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconStarFilled(
  { size = 24, stroke = 0, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="currentColor"
      stroke={color || "none"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M8.243 7.34l-6.38 .925l-.113 .023a1 1 0 0 0 -.44 1.684l4.622 4.499l-1.09 6.355l-.013 .11a1 1 0 0 0 1.464 .944l5.706 -3l5.693 3l.1 .046a1 1 0 0 0 1.352 -1.1l-1.091 -6.355l4.624 -4.5l.078 -.085a1 1 0 0 0 -.633 -1.62l-6.38 -.926l-2.852 -5.78a1 1 0 0 0 -1.794 0l-2.853 5.78z" />
    </svg>
  );
});
IconStarFilled.displayName = "IconStarFilled";

export const IconStrikethrough: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconStrikethrough(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 12l14 0" />
  <path d="M16 6.5a4 2 0 0 0 -4 -1.5h-1a3.5 3.5 0 0 0 0 7h2a3.5 3.5 0 0 1 0 7h-1.5a4 2 0 0 1 -4 -1.5" />
    </svg>
  );
});
IconStrikethrough.displayName = "IconStrikethrough";

export const IconSubtask: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconSubtask(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M6 9l6 0" />
  <path d="M4 5l4 0" />
  <path d="M6 5v11a1 1 0 0 0 1 1h5" />
  <path d="M12 8a1 1 0 0 1 1 -1h6a1 1 0 0 1 1 1v2a1 1 0 0 1 -1 1h-6a1 1 0 0 1 -1 -1l0 -2" />
  <path d="M12 16a1 1 0 0 1 1 -1h6a1 1 0 0 1 1 1v2a1 1 0 0 1 -1 1h-6a1 1 0 0 1 -1 -1l0 -2" />
    </svg>
  );
});
IconSubtask.displayName = "IconSubtask";

export const IconSun: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconSun(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M8 12a4 4 0 1 0 8 0a4 4 0 1 0 -8 0" />
  <path d="M3 12h1m8 -9v1m8 8h1m-9 8v1m-6.4 -15.4l.7 .7m12.1 -.7l-.7 .7m0 11.4l.7 .7m-12.1 -.7l-.7 .7" />
    </svg>
  );
});
IconSun.displayName = "IconSun";

export const IconTable: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconTable(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 5a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v14a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2v-14" />
  <path d="M3 10h18" />
  <path d="M10 3v18" />
    </svg>
  );
});
IconTable.displayName = "IconTable";

export const IconTag: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconTag(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M6.5 7.5a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M3 6v5.172a2 2 0 0 0 .586 1.414l7.71 7.71a2.41 2.41 0 0 0 3.408 0l5.592 -5.592a2.41 2.41 0 0 0 0 -3.408l-7.71 -7.71a2 2 0 0 0 -1.414 -.586h-5.172a3 3 0 0 0 -3 3" />
    </svg>
  );
});
IconTag.displayName = "IconTag";

export const IconTargetArrow: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconTargetArrow(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M11 12a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
  <path d="M12 7a5 5 0 1 0 5 5" />
  <path d="M13 3.055a9 9 0 1 0 7.941 7.945" />
  <path d="M15 6v3h3l3 -3h-3v-3l-3 3" />
  <path d="M15 9l-3 3" />
    </svg>
  );
});
IconTargetArrow.displayName = "IconTargetArrow";

export const IconTerminal2: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconTerminal2(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M8 9l3 3l-3 3" />
  <path d="M13 15l3 0" />
  <path d="M3 6a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2l0 -12" />
    </svg>
  );
});
IconTerminal2.displayName = "IconTerminal2";

export const IconThumbUp: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconThumbUp(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M7 11v8a1 1 0 0 1 -1 1h-2a1 1 0 0 1 -1 -1v-7a1 1 0 0 1 1 -1h3a4 4 0 0 0 4 -4v-1a2 2 0 0 1 4 0v5h3a2 2 0 0 1 2 2l-1 5a2 3 0 0 1 -2 2h-7a3 3 0 0 1 -3 -3" />
    </svg>
  );
});
IconThumbUp.displayName = "IconThumbUp";

export const IconTrash: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconTrash(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 7l16 0" />
  <path d="M10 11l0 6" />
  <path d="M14 11l0 6" />
  <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12" />
  <path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3" />
    </svg>
  );
});
IconTrash.displayName = "IconTrash";

export const IconTrendingUp: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconTrendingUp(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M3 17l6 -6l4 4l8 -8" />
  <path d="M14 7l7 0l0 7" />
    </svg>
  );
});
IconTrendingUp.displayName = "IconTrendingUp";

export const IconTruckDelivery: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconTruckDelivery(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 17a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" />
  <path d="M15 17a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" />
  <path d="M5 17h-2v-4m-1 -8h11v12m-4 0h6m4 0h2v-6h-8m0 -5h5l3 5" />
  <path d="M3 9l4 0" />
    </svg>
  );
});
IconTruckDelivery.displayName = "IconTruckDelivery";

export const IconUpload: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconUpload(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M4 17v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2 -2v-2" />
  <path d="M7 9l5 -5l5 5" />
  <path d="M12 4l0 12" />
    </svg>
  );
});
IconUpload.displayName = "IconUpload";

export const IconUser: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconUser(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M8 7a4 4 0 1 0 8 0a4 4 0 0 0 -8 0" />
  <path d="M6 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2" />
    </svg>
  );
});
IconUser.displayName = "IconUser";

export const IconUserCheck: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconUserCheck(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M8 7a4 4 0 1 0 8 0a4 4 0 0 0 -8 0" />
  <path d="M6 21v-2a4 4 0 0 1 4 -4h4" />
  <path d="M15 19l2 2l4 -4" />
    </svg>
  );
});
IconUserCheck.displayName = "IconUserCheck";

export const IconUserPlus: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconUserPlus(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M8 7a4 4 0 1 0 8 0a4 4 0 0 0 -8 0" />
  <path d="M16 19h6" />
  <path d="M19 16v6" />
  <path d="M6 21v-2a4 4 0 0 1 4 -4h4" />
    </svg>
  );
});
IconUserPlus.displayName = "IconUserPlus";

export const IconUsers: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconUsers(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 7a4 4 0 1 0 8 0a4 4 0 1 0 -8 0" />
  <path d="M3 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2" />
  <path d="M16 3.13a4 4 0 0 1 0 7.75" />
  <path d="M21 21v-2a4 4 0 0 0 -3 -3.85" />
    </svg>
  );
});
IconUsers.displayName = "IconUsers";

export const IconVariable: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconVariable(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M5 4c-2.5 5 -2.5 10 0 16m14 -16c2.5 5 2.5 10 0 16m-10 -11h1c1 0 1 1 2.016 3.527c.984 2.473 .984 3.473 1.984 3.473h1" />
  <path d="M8 16c1.5 0 3 -2 4 -3.5s2.5 -3.5 4 -3.5" />
    </svg>
  );
});
IconVariable.displayName = "IconVariable";

export const IconWallet: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconWallet(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M17 8v-3a1 1 0 0 0 -1 -1h-10a2 2 0 0 0 0 4h12a1 1 0 0 1 1 1v3m0 4v3a1 1 0 0 1 -1 1h-12a2 2 0 0 1 -2 -2v-12" />
  <path d="M20 12v4h-4a2 2 0 0 1 0 -4h4" />
    </svg>
  );
});
IconWallet.displayName = "IconWallet";

export const IconWaveSine: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconWaveSine(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M21 12h-2c-.894 0 -1.662 -.857 -1.761 -2c-.296 -3.45 -.749 -6 -2.749 -6s-2.5 3.582 -2.5 8s-.5 8 -2.5 8s-2.452 -2.547 -2.749 -6c-.1 -1.147 -.867 -2 -1.763 -2h-2" />
    </svg>
  );
});
IconWaveSine.displayName = "IconWaveSine";

export const IconX: Icon = forwardRef<SVGSVGElement, TablerIconProps>(function IconX(
  { size = 24, stroke = 2, color, className, ...restProps },
  forwardedRef
) {
  return (
    <svg
      ref={forwardedRef}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color || "currentColor"}
      strokeWidth={stroke}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...restProps}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
  <path d="M18 6l-12 12" />
  <path d="M6 6l12 12" />
    </svg>
  );
});
IconX.displayName = "IconX";

