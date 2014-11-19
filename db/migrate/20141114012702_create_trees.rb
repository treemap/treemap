class CreateTrees < ActiveRecord::Migration
  def change
    create_table :trees do |t|
      t.string :common_name
      t.string :latin_name # binomial nomenclature

      t.timestamps
    end
  end
end
