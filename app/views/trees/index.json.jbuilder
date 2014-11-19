json.array!(@trees) do |tree|
  json.extract! tree, :id
  json.url tree_url(tree, format: :json)
end
