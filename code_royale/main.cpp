#include <iostream>
#include <string>
#include <vector>
#include <algorithm>
#include <memory>
#include <cmath>

using namespace std;

class Unit {
public:
  Unit(int const& x, int const& y, int const& owner,
       int const& unitType, int const& health)
    :_x{x}, _y{y}, _owner{owner}, _type{unitType}, _health{health}
  {}

  Unit& operator=(Unit const& other) = default;

  // Getters and Setters
  int const& get_x() const {return _x;}
  int const& get_y() const {return _y;}
  int const& owner() const {return _owner;}
  int const& type() const {return _type;}
  int const& health() const {return _health;}

  // Other
  string stringify() {
    return "<Unit " + to_string(_owner)
      + " (" + to_string(_x) + ":" + to_string(_y) + ") health: "
      + to_string(_health) + ">";
  }

private:
  int _x;
  int _y;
  int _owner;
  int _type;
  int _health;
};

class Queen : public Unit {
public:
  Queen(int const& x, int const& y, int const& owner,
        int const& unitType, int const& health)
    :Unit(x, y, owner, unitType, health)
  {}
  Queen& operator=(Queen const& other) = default;
};

class Creep : public Unit {
public:
  Creep(int const& x, int const& y, int const& owner,
        int const& unitType, int const& health)
    :Unit(x, y, owner, unitType, health)
  {}
  bool is_knight() { return type() == 0;}
  bool is_archer() { return type() == 1;}
  bool is_giant() { return type() == 2;}
  bool is_mine() { return owner() == 0;}
};

class Site {
public:
  Site(int const& id, int const& x, int const& y, int const& radius)
    :_id{id}, _x{x}, _y{y}, _radius{radius}, _gold{-1}, _mine_size{-1},
     _type{-1}, _owner{-1}, _param1{-1}, _param2{-2}
  {}

  void update(int const& gold, int const& mine_size, int const& structureType,
              int const& owner, int const& param1, int const& param2) {
    _gold = gold;
    _mine_size = mine_size;
    _type = structureType;
    _owner = owner;
    _param1 = param1;
    _param2 = param2;
  }

  // Getters and Setters
  int const& id() const {return _id;}
  int const& get_x() const {return _x;}
  int const& get_y() const {return _y;}
  int const& get_radius() const {return _radius;}
  int const& get_param1() const {return _param1;}
  int const& get_param2() const {return _param2;}
  int const& get_gold() const {return _gold;}
  int const& get_mine_size() const {return _mine_size;}
  int const& type() const {return _type;}

  // Issers and Hassers
  bool is_mine() {return _owner == 0;}

  bool is_free() {return _type == -1;}
  bool is_a_mine() {return _type == 0;}
  bool is_tower() {return _type == 1;}
  bool is_barrack() {return _type == 2;}

  bool is_knight() {return _param2 == 0;}
  bool is_archer() {return _param2 == 1;}
  bool is_giant() {return _param2 == 2;}
  bool is_training() {return _param1 > 0;}

  // Others
  string const stringify() const {
    return "<Site " + to_string(id()) + ": (" + to_string(get_x()) + ":"
      + to_string(get_y()) + ") " + to_string(type()) + ">";
  }

private:
  // First part
  int _id;
  int _x;
  int _y;
  int _radius;
  // Second part
  int _gold;
  int _mine_size;
  int _type;
  int _owner;
  int _param1;
  int _param2;
};

class World {
public:
  World(int const& id)
    :_my_id{id}, _gold{-1},_touchedSite{-1}, _numUnits{-1}, _numSites{-1},
     _sites{}, _barracks{}, _creeps{}, _queen{nullptr}, _enemy_queen{nullptr},
     _target{nullptr}, _home_x{-1}, _home_y{-1}
  {}

  ~World() {
    for (Site* s: _sites) delete s;
    for (Creep* s: _creeps) delete s;
    delete _queen;
    delete _enemy_queen;
  }

  void add_site(Site* const& site) {
    _sites.push_back(site);
  }

  void add_queen(Queen* const& queen) {
    if (queen->owner() == _my_id)
      if (_queen != nullptr ) {
        *_queen = *queen;
      } else {
        _queen = queen;
        _home_x = queen->get_x();
        _home_y = queen->get_y();
      }
    else {
      if (_enemy_queen != nullptr) {
        *_enemy_queen = *queen;
      } else {
        _enemy_queen = queen;
      }
    }

  }

  void add_creep(Creep* const& creep) {
    _creeps.push_back(move(creep));
  }

  Site* site(int const& id) {
    for (Site* const& s : _sites) {
      if (s->id() == id) return s;
    }
    return nullptr;
  }

  void update_site(int const& id, int const& gold,
                   int const& mine_size, int const& structureType,
                   int const& owner, int const& param1, int const& param2) {
    site(id)->update(gold, mine_size, structureType, owner, param1, param2);
  }

  int distance(int const& x1, int const& y1,
               int const& x2, int const& y2) {
    int a{x1 - x2};
    int b{y1 - y2};
    return static_cast<int>(round(hypot(a,b)));
  }

  pair<Site*,int> closest_site() {
    Site* closest{nullptr};
    int closest_dist{99999};
    for (auto const& s : _sites) {
      int dist = distance(my_queen()->get_x(), my_queen()->get_y(),
                          s->get_x(), s->get_y()) - s->get_radius();
      if (dist < closest_dist) {
        closest = s;
        closest_dist = dist;
      }
    }
    return make_pair( closest, closest_dist );
  }

  pair<Site*, int> closest_free_site() {
    Site* closest{nullptr};
    int closest_dist{99999};
    for (auto const& s : _sites) {
      if (! s->is_free()) continue;
      int dist = distance(my_queen()->get_x(), my_queen()->get_y(),
                          s->get_x(), s->get_y()) - s->get_radius();
      if (dist < closest_dist) {
        closest = s;
        closest_dist = dist;
      }
    }
    return make_pair( closest, closest_dist );
  }

  pair<Site*, int> closest_free_safe_site() {
    vector<Site*> towers{};
    for (auto const& s : _sites) {
      if (s->is_free() || s->is_mine() || !s->is_tower()) continue;
      towers.push_back(s);
    }

    Site* closest{nullptr};
    int closest_dist{99999};
    for (auto const& s : _sites) {
      if (! s->is_free()) continue;
      bool safe{true};
      for (Site* t : towers) {
        int tower_dist = distance(t->get_x(), t->get_y(),
                                  s->get_x(), s->get_y());
        if ( (t->get_param2()) >= tower_dist ) {
          safe = false;
          break;
        }
      }
      if (!safe) continue;
      // tower safe free spot
      int dist = distance(my_queen()->get_x(), my_queen()->get_y(),
                          s->get_x(), s->get_y());
      if (dist < closest_dist) {
        closest = s;
        closest_dist = dist;
      }
    }

    return make_pair( closest, closest_dist );
  }

  pair<Site*, int> homiest_old_tower() {
    Site* homiest{nullptr};
    int homiest_dist{99999};
    int shortest_rad{99999};
    for (auto const& s : _sites) {
      if (!s->is_mine() || !s->is_tower()) continue;
      int dist = distance( _home_x, _home_y,
                           s->get_x(), s->get_y());
      dist -= s->get_radius();
      if (s->get_param2() < shortest_rad || dist < homiest_dist ) {
        homiest = s;
        homiest_dist = dist;
        shortest_rad = s->get_radius();
      }
    }
    return make_pair(homiest, homiest_dist);
  }

  pair<Site*, int> closest_old_tower() {
    Site* closest{nullptr};
    int closest_dist{99999};
    int shortest_rad{99999};
    for (auto const& s : _sites) {
      if (!s->is_mine() || !s->is_tower()) continue;
      int dist = distance(my_queen()->get_x(), my_queen()->get_y(),
                          s->get_x(), s->get_y());
      dist -= s->get_radius();
      if (dist < closest_dist || s->get_param2() < shortest_rad) {
        closest = s;
        closest_dist = dist;
        shortest_rad = s->get_radius();
      }
    }
    return make_pair(closest, closest_dist);
  }

  pair<Creep*, int> closest_archer() {
    Creep* closest{nullptr};

    int closest_dist{99999};
    for (auto const& c : _creeps) {
      if (!c->is_mine() || !c->is_archer()) continue;
      int dist = distance(my_queen()->get_x(), my_queen()->get_y(),
                          c->get_x(), c->get_y());
      if (dist < closest_dist) {
        closest = c;
        closest_dist = dist;
      }
    }
    return make_pair(closest, closest_dist);
  }

  pair<Site*, int> closest_enemy_tower() {
    Site* closest{nullptr};

    int closest_dist{99999};
    for (auto const& s : _sites) {
      if (s->is_mine() || !s->is_tower()) continue;
      int dist = distance(my_queen()->get_x(), my_queen()->get_y(),
                          s->get_x(), s->get_y());
      dist -= s->get_radius();
      if (dist < closest_dist) {
        closest = s;
        closest_dist = dist;
      }
    }
    return make_pair(closest, closest_dist);
  }

  pair<Site*,int> closest_enemy_barrack() {
    Site* closest{nullptr};

    int closest_dist{99999};
    for (auto const& s : _sites) {
      if (s->is_mine() || !s->is_barrack()) continue;
      int dist = distance(my_queen()->get_x(), my_queen()->get_y(),
                          s->get_x(), s->get_y());
      dist -= s->get_radius();
      if (dist < closest_dist) {
        closest = s;
        closest_dist = dist;
      }
    }
    return make_pair(closest, closest_dist);
  }

  pair<Site*,int> closest_mine() {
    Site* closest{nullptr};

    int closest_dist{99999};
    for (auto const& s : _sites) {
      if (!s->is_mine() || !s->is_a_mine()) continue;
      int dist = distance(my_queen()->get_x(), my_queen()->get_y(),
                          s->get_x(), s->get_y());
      dist -= s->get_radius();
      if (dist < closest_dist) {
        closest = s;
        closest_dist = dist;
      }
    }
    return make_pair(closest, closest_dist);
  }

  pair<Creep*, int> closest_enemy_knight() {
    Creep* closest{nullptr};
    int closest_dist{99999};
    for (Creep* c : _creeps) {
      if (!c->is_mine() && c->is_knight()) {
        auto dist = distance(_queen->get_x(), _queen->get_y(),
                             c->get_x(), c->get_y());
        if (dist <= closest_dist) {
          closest = c;
          closest_dist = dist;
        }
      }
    }
    return make_pair( closest, closest_dist );
  }

  int gpm() {
    int _gpm{0};
    for (Site* s: _sites) {
      if (s->is_mine() && s->is_a_mine()) _gpm += s->get_param1();
    }
    cerr << "GPM: " << _gpm << endl;
    return _gpm;
  }

  void clear_creeps() {
    for (Creep* c: _creeps) delete c;
    _creeps.clear();
  }

  // Constants
  int queen() {return -1;}
  int knight() {return 0;}
  int archer() {return 1;}
  int gpm_giant() {return 15;}

  // Getters and Setters
  int& gold() {return _gold;}
  vector<Site*> get_sites() {return _sites;}
  int& touched_site() {return _touchedSite;}
  int& numUnits() {return _numUnits;}
  int& numSites() {return _numSites;}
  Queen* my_queen() const {return _queen;}
  Site*& target() {return _target;}
  int& home_x() {return _home_x;}
  int& home_y() {return _home_y;}

  // Issers and Hassers
  bool touches_site() {return _touchedSite != -1;}

  bool has_archers(int const& n) const {
    int found{0};
    for (Creep* c : _creeps) {
      if (c->is_mine() && c->is_archer()) ++found;
    }
    cerr << "Archers found: " << found << endl;
    return n <= found;
  }

  bool has_giants(int const& n) const {
    int found{0};
    for (Creep* c : _creeps) {
      if (c->is_mine() && c->is_giant()) ++found;
    }
    cerr << "Giants found: " << found << endl;
    return n <= found;
  }

  bool has_knights(int const& n) const {
    int found{0};
    for (Creep* c : _creeps) {
      if (c->is_mine() && c->is_knight()) ++found;
    }
    cerr << "Knights found: " << found << endl;
    return n <= found;
  }

  bool has_mines(int const& n) const {
    int found{0};
    for (Site* s : _sites) {
      if (s->is_mine() && s->is_a_mine()) ++found;
    }
    cerr << "Mines found: " << found << endl;
    return n <= found;
  }

  bool has_towers(int const& n) const {
    int found{0};
    for (Site* s : _sites) {
      if (s->is_mine() && s->is_tower()) ++found;
    }
    cerr << "Towers found: " << found << endl;
    return n <= found;
  }

  bool has_enemy_towers(int const& n = 1) const {
    int found{0};
    for (Site* s : _sites) {
      if (!s->is_mine() && s->is_tower()) ++found;
    }
    cerr << "Enemy Towers found: " << found << endl;
    return n <= found;
  }

  bool has_giant_barracks(int const& n) const {
    int found{0};
    for (Site* s : _sites) {
      if (s->is_mine() && s->is_barrack() && s->is_giant() ) ++found;
    }
    cerr << "Barracks Giant found: " << found << endl;
    return n <= found;
  }

  bool has_archer_barracks(int const& n) const {
    int found{0};
    for (Site* s : _sites) {
      if (s->is_mine() && s->is_barrack() && s->is_archer() ) ++found;
    }
    cerr << "Barracks Archer found: " << found << endl;
    return n <= found;
  }

  bool has_knight_barracks(int const& n) const {
    int found{0};
    for (Site* s : _sites) {
      if (s->is_mine() && s->is_barrack() && s->is_knight() ) ++found;
    }
    cerr << "Barracks Knight found: " << found << endl;
    return n <= found;
  }

  bool knights_are_too_close() {
    int limit = my_queen()->health() < 30 ? 1000 : 600;
    for (Creep* c : _creeps) {
      if (!c->is_mine() && c->is_knight()) {
        auto d = distance(_queen->get_x(), _queen->get_y(),
                          c->get_x(), c->get_y());
        if (d <= limit)
          return true;
      }
    }
    return false;
  }

  bool tower_too_close() {
    for (Site* s: _sites) {
      if (s->is_mine() || !s->is_tower()) continue;
      auto dist = distance(_queen->get_x(), _queen->get_y(),
                           s->get_x(), s->get_y());
      if (dist <= s->get_param2() + 60) // Some extra for queen
        return true;
    }
    return false;
  }
  bool enemy_health_low() {
    return _enemy_queen->health() < 10;
  }
  // Others
  void p_sites() const {
    for (Site* const& s : _sites) cerr << s->stringify() << endl;
  }

private:
  int _my_id;
  int _gold;
  int _touchedSite;
  int _numUnits;
  int _numSites;
  vector<Site*> _sites;
  vector<Site*> _barracks;
  vector<Creep*> _creeps;
  Queen* _queen;
  Queen* _enemy_queen;
  Site* _target;
  int _home_x;
  int _home_y;
};

/**********************************************************
 * MAIN
 **********************************************************/

int main()
{
  World world{0};
  cin >> world.numSites(); cin.ignore();
  for (int i = 0; i < world.numSites(); i++) {
    int siteId;
    int x;
    int y;
    int radius;
    cin >> siteId >> x >> y >> radius; cin.ignore();
    world.add_site(new Site{siteId,x,y,radius});
  }

  // game loop
  while (1) {
    // touchedSite -1 if none
    vector<string> commands{2,""};
    cin >> world.gold() >> world.touched_site(); cin.ignore();
    for (int i = 0; i < world.numSites(); i++) {
      int siteId;
      int gold; // used in future leagues
      int mine_size; // used in future leagues
      int structureType; // -1 = No structure, 2 = Barracks
      int owner; // -1 = No structure, 0 = Friendly, 1 = Enemy
      int param1;
      int param2;
      cin >> siteId >> gold >> mine_size
          >> structureType >> owner >> param1 >> param2; cin.ignore();

      world.update_site(siteId, gold, mine_size,
                        structureType, owner, param1, param2);
    }
    cin >> world.numUnits(); cin.ignore();
    world.clear_creeps();
    for (int i = 0; i < world.numUnits(); i++) {
      int x;
      int y;
      int owner;
      int unitType;
      int health;
      cin >> x >> y >> owner >> unitType >> health; cin.ignore();
      if (unitType == world.queen()) {
        world.add_queen(new Queen(x,y,owner,unitType,health));
      } else {
        world.add_creep(new Creep(x,y,owner,unitType,health));
      }
    }

    /*
     *
     * GAME LOGIC STARTS HERE :D
     *
     */

    Site* touched_site = world.site(world.touched_site());
    int gpm{world.gpm()};
    if (world.touches_site() && touched_site->is_free()) {
      // At a site with no building on it
      string next_building{};
      if (world.knights_are_too_close()) {
        next_building = "TOWER";
      } else if (touched_site->get_gold() > 0 && !world.has_mines(2)) {
        next_building = "MINE";
      } else if (gpm > world.gpm_giant() && world.gold() > 300 && !world.has_giant_barracks(1)) {
        next_building = "BARRACKS-GIANT";
      } else if (gpm > 4 && !world.has_knight_barracks(2)) {
        next_building = "BARRACKS-KNIGHT";
      } else if (gpm > 5 && !world.has_archer_barracks(1)) {
        next_building = "BARRACKS-ARCHER";
      } else if (gpm > 2 && !world.has_knight_barracks(1)) {
        next_building = "BARRACKS-KNIGHT";
      } else if (touched_site->get_gold() > 0){
        next_building = "MINE";
      } else {
        next_building = "TOWER";
      }
      commands.at(0) = "BUILD " + to_string(touched_site->id()) + " " + next_building;
    } else {
      pair<Site*, int> closest_tower = world.closest_old_tower();
      pair<Site*, int> homiest_tower = world.homiest_old_tower();
      pair<Creep*, int> closest_archer = world.closest_archer();
      pair<Site*, int> closest_free = (world.has_enemy_towers(1) ? world.closest_free_safe_site() : world.closest_free_site());
      pair<Site*, int> closest_enemy_barrack = world.closest_enemy_barrack();
      //pair<Creep*, int> closest_enemy_knight = world.closest_enemy_knight();

      if (world.knights_are_too_close()) {
        if (homiest_tower.first != nullptr && world.has_towers(3) && homiest_tower.first->get_param1() < 500) {
          cerr << "Repairing homiest oldest tower: " << closest_tower.first->stringify() << endl;

          commands.at(0) = "BUILD " + to_string(homiest_tower.first->id()) + " TOWER";
        } else if (closest_tower.first != nullptr && world.has_towers(3)) {
          cerr << "Repairing closest oldest tower: " << closest_tower.first->stringify() << endl;

          commands.at(0) = "BUILD " + to_string(closest_tower.first->id()) + " TOWER";
        } else if (closest_free.first != nullptr) {
          cerr << "Building a tower at " << closest_free.first->stringify() << endl;

          commands.at(0) = "BUILD " + to_string(closest_free.first->id()) + " TOWER";
        } else if (world.touches_site() && (
                                            (touched_site->is_barrack() && !touched_site->is_training())
                                            || touched_site->is_mine() || touched_site->is_tower() )){
          cerr << "Building a tower at " << closest_free.first->stringify() << endl;

          commands.at(0) = "BUILD " + to_string(closest_free.first->id()) + " TOWER";
        } else if (world.has_archers(1)) {
          cerr << "Moving to closest Archer " << closest_archer.first->stringify() << endl;

          commands.at(0) = "MOVE " + to_string(closest_archer.first->get_x())
            + " " + to_string(closest_archer.first->get_y());
        } else {
          cerr << "RUNNING \"HOME\" :(" << endl;

          commands.at(0) = "MOVE 0 0";
        }
      } else if (closest_enemy_barrack.second < 300 && world.my_queen()->health() > 40) {
        cerr << "Destroying Barrack " << closest_enemy_barrack.first->stringify() << endl;
        commands.at(0) = "MOVE " + to_string(closest_enemy_barrack.first->get_x()) + " " + to_string(closest_enemy_barrack.first->get_y());

      } else if (world.touches_site() && touched_site->is_tower() && touched_site->get_param1() < 550 ) {
        cerr << "Upgrading Tower " << touched_site->stringify() << endl;
        commands.at(0) = "BUILD " + to_string(touched_site->id()) + " TOWER";

      } else if (world.touches_site() && touched_site->is_a_mine() && touched_site->get_mine_size() > 2 && touched_site->get_param1() < 3) {
        cerr << "Upgrading Mine" << touched_site->stringify() << endl;

        commands.at(0) = "BUILD " + to_string(touched_site->id()) + " MINE";

      } else if (!world.has_knight_barracks(1) && world.touches_site()
                 && touched_site->is_a_mine() && touched_site->get_mine_size() == 1) {
        cerr << "Replacing Mine" << touched_site->stringify() << endl;
        commands.at(0) = "BUILD " + to_string(touched_site->id()) + " BARRACKS-KNIGHT";

      } else if (!world.has_archer_barracks(1) && world.touches_site()
                 && touched_site->is_a_mine() && touched_site->get_mine_size() == 1) {
        cerr << "Replacing Mine" << touched_site->stringify() << endl;
        commands.at(0) = "BUILD " + to_string(touched_site->id()) + " BARRACKS-ARCHER";

      } else if (!world.has_giant_barracks(1) && world.touches_site()
                 && touched_site->is_a_mine() && touched_site->get_mine_size() == 1) {
        cerr << "Replacing Mine" << touched_site->stringify() << endl;
        commands.at(0) = "BUILD " + to_string(touched_site->id()) + " BARRACKS-GIANT";
      } else {
        if (closest_free.first != nullptr) {
          cerr << "Moving to closest site: " << closest_free.first->stringify() << endl;

          commands.at(0) = "MOVE " + to_string(closest_free.first->get_x())
            + " " + to_string(closest_free.first->get_y());
        } else if (closest_enemy_barrack.first != nullptr) {
          cerr << "Moving to closest enemy barracks: " << closest_enemy_barrack.first->stringify() << endl;

          commands.at(0) = "MOVE " + to_string(closest_enemy_barrack.first->get_x())
            + " " + to_string(closest_enemy_barrack.first->get_y());
        } else {
          cerr << "NOTHING TO DO" << endl;
        }
      }
    }

    /*
      FOR TRAINING
    */
    commands.at(1) = "TRAIN";
    for (Site* s: world.get_sites()) {
      if (s->is_mine() && !s->is_free() && s->is_barrack()) {
        // if we have one giant out, don't make more.
        if (s->is_training()) continue;
        else if (s->is_knight() && world.enemy_health_low()) {
          commands.at(1) += " " + to_string(s->id());
          break;
        }
        else if (s->is_knight() && !world.has_knights(12))
          commands.at(1) += " " + to_string(s->id());
        else if(s->is_archer() && !world.has_archers(3))
          commands.at(1) += " " + to_string(s->id());
        else if (gpm > world.gpm_giant() && s->is_giant() && !world.has_giants(1))
          commands.at(1) += " " + to_string(s->id());
      }
    }

    // First line: A valid queen action
    // Second line: A set of training instructions
    for (string c : commands) cout << c << endl;
  }
}
