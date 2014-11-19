# From the data collected one can calculate the number of trees per acre, the basal area, the volume of trees in an area, and the value of the timber.
class TreeInventory < ActiveRecord::Base
  belongs_to :property

  has_many :tree_diameters

  enum unit: [
         :imperial,
         :metric,
       ]

  def mean_tree_diameters
    tree_map = Hash.new { [] }

    self.tree_diameters.each do |td|
      tree_map[td.tree] << td.diameter
    end

    tree_means = Hash.new { 0 }
    tree_map.each do |k, v|
      tree_means[k] = v.reduce(:+) / tree_means[k].length rescue 0
    end

    # Take an average by tree.

    tree_means
  end


  # http://en.wikipedia.org/wiki/Basal_area
  #
  # The basal area of a forest stand can be found by adding the basal
  # areas (as calculated above) of all of the trees in an area and
  # dividing by the area of land in which the trees were
  # measured. Basal area is generally expressed as ft^2/acre or m^2/ha.

  def basal_area
    basal_sum = self.tree_diameters.map(&:basal_area).reduce(:+)

    basal_sum / property.area
  end
end
