import classNames from "classnames";

export type AvatarSize = "6" | "8" | "10" | "12" | "14";
const sizeClassMapping: Record<AvatarSize, string> = {
  "6": "h-6 w-6",
  "8": "h-8 w-8",
  "10": "h-10 w-10",
  "12": "h-12 w-12",
  "14": "h-14 w-14",
};
export default function UserAvatar({
  size = "10",
  src,
  alt,
}: {
  size?: AvatarSize;
  src?: string;
  alt?: string;
}) {
  if (src) {
    return (
      <img
        className={classNames("rounded-full", sizeClassMapping[size])}
        src={src}
        alt={alt}
      />
    );
  }
  return (
    <span
      className={classNames(
        "inline-block overflow-hidden rounded-full bg-gray-100",
        sizeClassMapping[size]
      )}
    >
      <svg
        className="h-full w-full text-gray-300"
        fill="currentColor"
        viewBox="0 0 24 24"
      >
        <path d="M24 20.993V24H0v-2.996A14.977 14.977 0 0112.004 15c4.904 0 9.26 2.354 11.996 5.993zM16.002 8.999a4 4 0 11-8 0 4 4 0 018 0z" />
      </svg>
    </span>
  );
}
