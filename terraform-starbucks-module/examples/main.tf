terraform {
  required_providers {
    starbucks = {
      source  = "registry.terraform.io/starbucks/starbucks"
      version = "~> 0.1"
    }
  }
}

provider "starbucks" {
  api_key = var.api_key
  region  = var.region
  endpoint = "https://api.starbucks.com/v1"
}

variable "api_key" {
  description = "Starbucks API Key"
  type        = string
  sensitive   = true
}

# Create multiple stores
resource "starbucks_store" "flagship_stores" {
  for_each = {
    seattle = {
      name         = "Seattle Flagship"
      store_number = "10001"
      address      = "2401 Utah Ave S"
      city         = "Seattle"
      state        = "WA"
      zip_code     = "98134"
      latitude     = 47.5759
      longitude    = -122.3263
      has_drive_thru = false
      store_type   = "reserve"
      capacity     = 150
    }
    new_york = {
      name         = "New York Reserve Roastery"
      store_number = "10002"
      address      = "61 9th Ave"
      city         = "New York"
      state        = "NY"
      zip_code     = "10011"
      latitude     = 40.7423
      longitude    = -74.0060
      has_drive_thru = false
      store_type   = "reserve"
      capacity     = 200
    }
  }

  name           = each.value.name
  store_number   = each.value.store_number
  address        = each.value.address
  city           = each.value.city
  state          = each.value.state
  zip_code       = each.value.zip_code
  country        = "US"
  phone_number   = "+1-800-STARBUC"
  
  latitude       = each.value.latitude
  longitude      = each.value.longitude
  
  opening_hours  = "Mon-Sun: 6AM-10PM"
  
  has_drive_thru  = each.value.has_drive_thru
  has_wifi        = true
  has_mobile_order = true
  
  capacity       = each.value.capacity
  store_type     = each.value.store_type
  manager_email  = "manager.${each.key}@starbucks.com"
}

# Create employees for Seattle store
resource "starbucks_employee" "seattle_team" {
  for_each = {
    manager = {
      first_name = "John"
      last_name  = "Smith"
      position   = "store_manager"
      rate       = 28.50
      is_supervisor = true
    }
    supervisor = {
      first_name = "Sarah"
      last_name  = "Johnson"
      position   = "shift_supervisor"
      rate       = 22.00
      is_supervisor = true
    }
    barista1 = {
      first_name = "Mike"
      last_name  = "Davis"
      position   = "barista"
      rate       = 18.00
      is_supervisor = false
    }
    barista2 = {
      first_name = "Emily"
      last_name  = "Wilson"
      position   = "barista"
      rate       = 18.50
      is_supervisor = false
    }
  }

  employee_number = "EMP-SEA-${upper(each.key)}"
  first_name      = each.value.first_name
  last_name       = each.value.last_name
  email           = "${lower(each.value.first_name)}.${lower(each.value.last_name)}@starbucks.com"
  phone_number    = "+1-206-555-${format("%04d", index(keys(starbucks_employee.seattle_team), each.key))}"
  
  store_id        = starbucks_store.flagship_stores["seattle"].id
  position        = each.value.position
  hire_date       = "2024-01-01"
  
  hourly_rate     = each.value.rate
  
  is_barista          = true
  is_shift_supervisor = each.value.is_supervisor
  is_certified        = true
  
  employment_type = "full_time"
}

# Create menu items
resource "starbucks_menu_item" "signature_drinks" {
  for_each = {
    psl = {
      name     = "Pumpkin Spice Latte"
      price    = 5.95
      calories = 380
      seasonal = true
    }
    caramel_macchiato = {
      name     = "Caramel Macchiato"
      price    = 5.45
      calories = 250
      seasonal = false
    }
    cold_brew = {
      name     = "Cold Brew"
      price    = 4.95
      calories = 5
      seasonal = false
    }
  }

  name        = each.value.name
  category    = "coffee"
  size        = "grande"
  price       = each.value.price
  calories    = each.value.calories
  
  is_available = true
  is_seasonal  = each.value.seasonal
}

# Create inventory for Seattle store
resource "starbucks_inventory" "seattle_inventory" {
  for_each = {
    beans = {
      item     = "Pike Place Roast"
      type     = "beans"
      quantity = 100
      unit     = "lbs"
      reorder  = 20
    }
    milk = {
      item     = "Whole Milk"
      type     = "milk"
      quantity = 50
      unit     = "gallons"
      reorder  = 10
    }
    cups = {
      item     = "Grande Cups"
      type     = "cups"
      quantity = 5000
      unit     = "count"
      reorder  = 1000
    }
  }

  store_id       = starbucks_store.flagship_stores["seattle"].id
  item_name      = each.value.item
  item_type      = each.value.type
  quantity       = each.value.quantity
  unit           = each.value.unit
  reorder_level  = each.value.reorder
  last_restocked = "2024-10-01"
}

# Create promotions
resource "starbucks_promotion" "seasonal_promos" {
  for_each = {
    fall = {
      name        = "Fall Favorites"
      description = "Enjoy 20% off all fall drinks"
      discount    = 20.0
      start       = "2024-09-01"
      end         = "2024-11-30"
      code        = "FALL2024"
    }
    rewards = {
      name        = "Rewards Member Special"
      description = "Extra star rewards for members"
      discount    = 15.0
      start       = "2024-10-01"
      end         = "2024-12-31"
      code        = "REWARDS15"
    }
  }

  name                = each.value.name
  description         = each.value.description
  discount_percentage = each.value.discount
  start_date          = each.value.start
  end_date            = each.value.end
  applies_to          = "all"
  is_active           = true
  promo_code          = each.value.code
}

# Data source examples
data "starbucks_store" "lookup_seattle" {
  id = starbucks_store.flagship_stores["seattle"].id
}

data "starbucks_stores" "washington_stores" {
  state = "WA"
}

# Outputs
output "seattle_store_id" {
  value = starbucks_store.flagship_stores["seattle"].id
}

output "total_employees" {
  value = length(starbucks_employee.seattle_team)
}

output "washington_store_count" {
  value = length(data.starbucks_stores.washington_stores.stores)
}

output "promotion_codes" {
  value = {
    for k, v in starbucks_promotion.seasonal_promos : k => v.promo_code
  }
}
