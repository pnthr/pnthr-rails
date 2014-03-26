class InstancesController < ApplicationController
  before_action :set_instance, only: [:show, :edit, :update, :destroy]

  # GET /applications
  def index
    @instances = Instance.all
  end

  # GET /instances/1
  def show
  end

  # GET /instances/new
  def new
    @instance = Instance.new
  end

  # GET /instances/1/edit
  def edit
  end

  # POST /instances
  def create
    @instance = Instance.new(instance_params)

    @instance.identifier = Digest::MD5.hexdigest(nil, false)
    @instance.secret = Digest::MD5.hexdigest(nil, false)
    @instance.password = Digest::MD5.hexdigest(instance_params[:password])

    if @instance.save
      redirect_to @instance, notice: 'Instance was successfully created.'
    else
      render action: 'new'
    end
  end

  # PATCH/PUT /instances/1
  def update
    if @instance.update(instance_params)
      redirect_to @instance, notice: 'Instance was successfully updated.'
    else
      render action: 'edit'
    end
  end

  # DELETE /instances/1
  def destroy
    @instance.destroy
    redirect_to instances_url, notice: 'Instance was successfully destroyed.'
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_instance
      @instance = Instance.find(params[:id])
    end

    # Only allow a trusted parameter "white list" through.
    def instance_params
      params.require(:instance).permit(:name, :description, :identifier, :secret, :password)
    end
end
