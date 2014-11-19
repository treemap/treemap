
mkdir zip_files
cd zip_files
# cat ../tree_files.csv | awk '{print $1}' | sed -e "s/,//g" -e "s/^/http:\/\/esp.cr.usgs.gov\/data\/little\//g" | xargs wget
# ls *.zip | parallel unzip -a

# for shp in `ls *.shp`; do
#     file_head=`echo $shp | sed -e "s/.shp//g"`
#     shp2pgsql -s 4326 -I -c -W UTF-8 $shp $file_head > $file_head.sql
# done

for shp in `ls *.sql`; do
    psql forestly_development < $shp
done
