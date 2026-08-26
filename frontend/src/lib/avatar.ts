export const getAvatarUrl = (username: string) => {
  const safeName = username || "User";
  return `https://api.dicebear.com/7.x/initials/svg?seed=${encodeURIComponent(safeName)}&radius=50&fontFamily=sans-serif&fontWeight=700`;
};
