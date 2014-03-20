class InstancesController < ApplicationController
  before_action :set_instance, only: [:show, :edit, :update, :destroy]

  # GET /applications
  def index
    @applications = Instance.all
  end

  # GET /applications/1
  def show
  end

  # GET /applications/new
  def new
    @application = Instance.new
  end

  # GET /applications/1/edit
  def edit
  end

  # POST /applications
  def create
    @application = Instance.new(instance_params)

    @application.identifier = SecureRandom.urlsafe_base64(nil, false)
    @application.secret = SecureRandom.urlsafe_base64(nil, false)
    @application.password = Digest::MD5.hexdigest(instance_params[:password])

    if @application.save
      redirect_to @application, notice: 'Instance was successfully created.'
    else
      render action: 'new'
    end
  end

  # PATCH/PUT /applications/1
  def update
    if @application.update(instance_params)
      redirect_to @application, notice: 'Instance was successfully updated.'
    else
      render action: 'edit'
    end
  end

  # DELETE /applications/1
  def destroy
    @application.destroy
    redirect_to applications_url, notice: 'Instance was successfully destroyed.'
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_instance
      @application = Instance.find(params[:id])
    end

    # Only allow a trusted parameter "white list" through.
    def instance_params
      params.require(:instance).permit(:name, :description, :identifier, :secret, :password)
    end
end
