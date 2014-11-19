library(rworldmap)

coords <- read.csv("coords.csv", header = FALSE)
colnames(coords) <- c("lon", "lat")

coords2 <- read.csv("coords.csv", header = FALSE)
colnames(coords2) <- c("lon", "lat")


newmap <- getMap(resolution = "high")
plot(newmap, asp = 1)
# points(coords2$lon, coords2$lat, col = "red", cex = 1)
points(coords$lon, coords$lat, col = "red", cex = 1)
