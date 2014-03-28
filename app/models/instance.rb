class Instance
  include Mongoid::Document
  include Mongoid::Timestamps

  field :name, type: String
  field :description, type: String
  field :secret, type: String
  field :password, type: String

  validates_presence_of :name, :secret, :password
end
