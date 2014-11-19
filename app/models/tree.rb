class Tree < ActiveRecord::Base

  def geodata
    tree_geoms = TreeGeom.where(latin_name: self.latin_name)
    tree_geoms.map do |i|
      JSON.parse(TreeGeom.connection.execute("select ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geodata from tree_geoms where id = #{i.id}").first["geodata"])["coordinates"].first.first rescue []
    end
  end

end
