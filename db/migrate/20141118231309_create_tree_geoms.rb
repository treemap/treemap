class CreateTreeGeoms < ActiveRecord::Migration
  def change
    create_table :tree_geoms do |t|
      t.string :latin_name

      t.multi_polygon :geom, :srid => 4326

      t.timestamps
    end
  end
end
